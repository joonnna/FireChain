package core

import (
	"math/rand"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/joonnna/blocks/protobuf"
	"github.com/joonnna/ifrit/logger"
)

type state struct {
	log *logger.Log

	pool *entryPool

	sync.RWMutex
	peerMap    map[string]*peer
	localPeer  *peer
	inProgress *block
}

func newState(localId string, log *logger.Log) *state {
	peerMap := make(map[string]*peer)
	localPeer := &peer{id: localId}
	peerMap[localId] = localPeer

	return &state{
		peerMap:    peerMap,
		pool:       newEntryPool(),
		localPeer:  localPeer,
		log:        log,
		inProgress: createBlock(nil),
	}
}

func (s *state) bytes() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	m := &blockchain.State{
		Peers:          make(map[string]*blockchain.PeerState),
		PendingEntries: make(map[string]bool),
	}
	for _, p := range s.peerMap {
		if entries := p.getEntries(); entries != nil {
			s := &blockchain.PeerState{
				Id:          p.id,
				Epoch:       p.getEpoch(),
				EntryHashes: entries,
				RootHash:    p.getRootHash(),
				PrevHash:    p.getPrevHash(),
			}

			m.Peers[p.id] = s
		}
	}

	for _, e := range s.pool.getAllPending() {
		key := string(e.hash)
		if _, exists := m.PendingEntries[key]; !exists {
			m.PendingEntries[key] = true
		}
	}

	//s.log.Debug.Println("Number of pending entries: ", len(m.GetPendingEntries()))

	return proto.Marshal(m)
}

func (s *state) merge(other *blockchain.State) ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	resp := &blockchain.StateResponse{}

	if peers := other.GetPeers(); peers != nil {
		s.mergePeers(peers)
	} else {
		s.log.Debug.Println("No peer state")
	}

	if entries := other.GetPendingEntries(); entries != nil {
		resp.MissingEntries = s.pool.diff(entries)
	} else {
		s.log.Debug.Println("No pending entries")
	}

	b, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *state) mergePeers(peers map[string]*blockchain.PeerState) {
	var favourites []*peer
	var numVotes uint = 0
	blockVotes := make(map[string]uint)

	for id, p := range peers {
		if id == s.localPeer.id {
			continue
		}
		s.mergePeer(p)
	}

	for _, p := range s.peerMap {
		key := string(p.getRootHash())

		if _, exists := blockVotes[key]; !exists {
			blockVotes[key] = 1
			if numVotes == 0 {
				numVotes = 1
				favourites = append(favourites, p)
			}
		} else {
			blockVotes[key]++
			if votes := blockVotes[key]; votes > numVotes {
				favourites = nil
				favourites = append(favourites, p)
				numVotes = votes
			} else if votes == numVotes {
				favourites = append(favourites, p)
			}
		}
	}

	s.log.Debug.Println("Block options: ", len(blockVotes))
	s.log.Debug.Println("Number of candidates: ", len(favourites))
	if numFavs := len(favourites); numFavs >= 1 {
		idx := 0
		if numFavs > 1 {
			idx = rand.Int() % numFavs
		}

		favPeer := favourites[idx]
		s.localPeer.increment(favPeer.getRootHash(), favPeer.getPrevHash(), favPeer.getEntries())
		s.log.Debug.Println("Favourite rootHash: ", string(favPeer.getRootHash()))
		s.log.Debug.Println("Favourite prevHash: ", string(favPeer.getPrevHash()))
	}
}

func (s *state) mergePeer(p *blockchain.PeerState) {
	if p.GetEntryHashes() == nil {
		return
	}
	if existingPeer := s.peerMap[p.GetId()]; existingPeer != nil {
		if p.GetEpoch() > existingPeer.getEpoch() {
			existingPeer.update(p.GetEpoch(), p.GetEntryHashes(), p.GetRootHash(), p.GetPrevHash())
		}
	} else {
		newPeer := &peer{
			id:       p.GetId(),
			entries:  p.GetEntryHashes(),
			rootHash: p.GetRootHash(),
			prevHash: p.GetPrevHash(),
			epoch:    p.GetEpoch(),
		}

		s.peerMap[newPeer.id] = newPeer
	}
}

func (s *state) add(e *entry) {
	s.Lock()
	defer s.Unlock()

	s.pool.addPending(e)

	err := s.inProgress.add(e)
	if err == errFullBlock && !s.localPeer.haveFavourite() {
		s.localPeer.addBlock(s.inProgress)
	}
	/*
		else if err != nil {
			s.log.Err.Println(err)
		}
	*/
}

func (s *state) exists(key string) bool {
	s.RLock()
	defer s.RUnlock()

	return s.pool.exists(key)
}

func (s *state) newRound(prevHash []byte) *block {
	s.Lock()
	defer s.Unlock()

	s.log.Info.Println("Picking favourite block")

	new := createBlock(prevHash)

	if s.localPeer.getEntries() == nil {
		s.log.Debug.Println("no entries, the fuck")
	}

	for _, e := range s.localPeer.getEntries() {
		key := string(e)
		ent := s.pool.getPending(key)
		if ent == nil {
			s.log.Err.Println("Missing entry in fav block")
			continue
		}

		err := new.add(ent)
		if err != nil {
			s.log.Err.Println(err)
		} else {
			s.pool.removePending(key)
			s.pool.addConfirmed(ent)
		}
	}

	if eq := new.cmpRootHash(s.localPeer.getRootHash()); !eq {
		s.log.Err.Println("Have different root hash for new fav block")
	}

	for _, p := range s.peerMap {
		p.reset()
	}

	s.inProgress = createBlock(new.getRootHash())

	if full := s.pool.fillWithPending(s.inProgress); full {
		s.localPeer.addBlock(s.inProgress)
	}

	return new
}

/*
func (s *state) getPeerMap() []*peer {
	c.peerMapMutex.RLock()
	defer c.peerMapMutex.RUnlock()

	ret := make([]*peer, len(c.peerMap))
	idx := 0

	for _, p := range c.peerMap {
		ret[idx] = p
		idx++
	}

	return ret
}

// must hold write access mutex before calling
func (s *state) addPeer(p *peer) {
	if _, exists := c.peerMap[p.id]; exists {
		return
	}

	c.peerMap[p.id] = p
}

// must hold reac access mutex before calling
func (s *state) getPeer(key string) *peer {
	if p, exists := c.peerMap[key]; !exists {
		return nil
	} else {
		return p
	}
}


func (s *state) reset() {
	defer s.Unlock()

	for _, p := range s.peerMap {
		p.reset()
	}
}

func (s *state) freeze() {
	s.Lock()
}
*/
