package core

import (
	"bytes"
	"math/rand"
	"sync"

	"github.com/golang/protobuf/proto"

	log "github.com/inconshreveable/log15"
	"github.com/joonnna/blocks/protobuf"
)

type state struct {
	sync.RWMutex
	peerMap    map[string]*peer
	pool       *entryPool
	localPeer  *peer
	inProgress *block

	// Experiments only
	round                  func() uint32
	prevRoundNumber        uint32
	convergeMap            map[string]uint32
	converged              bool
	experimentParticipants int
}

func newState(localId string, httpAddr string, round func() uint32, hosts int) *state {
	peerMap := make(map[string]*peer)
	localPeer := &peer{
		id:       localId,
		httpAddr: httpAddr,
	}
	peerMap[localId] = localPeer

	return &state{
		peerMap:    peerMap,
		pool:       newEntryPool(),
		localPeer:  localPeer,
		inProgress: createBlock(nil),
		round:      round,
		// Experiment hosts might fail, will never register convergence if anyone fails
		// 10% failsafe
		experimentParticipants: int(float32(hosts) * 0.90),
		convergeMap:            make(map[string]uint32),
	}
}

func (s *state) bytes() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	m := &blockchain.State{
		Peers:          make(map[string]*blockchain.PeerState),
		MissingEntries: s.pool.getAllMissing(),
	}
	for _, p := range s.peerMap {
		if entries := p.getEntries(); entries != nil {
			s := &blockchain.PeerState{
				Id:          p.id,
				Epoch:       p.getEpoch(),
				EntryHashes: entries,
				RootHash:    p.getRootHash(),
				PrevHash:    p.getPrevHash(),
				HttpAddr:    p.httpAddr,
			}

			m.Peers[p.id] = s
		}
	}

	return proto.Marshal(m)
}

func (s *state) merge(other *blockchain.State) ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	if peers := other.GetPeers(); peers != nil {
		s.mergePeers(peers)
	} else {
		log.Debug("No peer state")
	}

	if entries := other.GetMissingEntries(); entries != nil {
		resp := &blockchain.StateResponse{
			Entries: s.pool.diff(entries),
		}

		log.Debug("Request for missing entries", "amount", len(entries))

		b, err := proto.Marshal(resp)
		if err != nil {
			return nil, err
		}

		return b, nil
	} else {
		log.Debug("Returning nil response")
		return nil, nil
	}
}

func (s *state) mergePeers(peers map[string]*blockchain.PeerState) {
	var favourites []*peer
	var numVotes int = 0
	blockVotes := make(map[string]int)

	for id, p := range peers {
		if id == s.localPeer.id {
			continue
		}
		s.mergePeer(p)
	}

	for _, p := range s.peerMap {
		// We do not consider peers with different previous blocks,
		// risk of creating a fork.
		if eq := bytes.Equal(s.localPeer.getPrevHash(), p.getPrevHash()); !eq {
			continue
		}
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

	//log.Debug("Block voting", "Options", len(blockVotes), "Favourites", len(favourites))

	if numFavs := len(favourites); numFavs >= 1 {
		idx := 0
		if numFavs > 1 {
			idx = rand.Int() % numFavs
		}

		favPeer := favourites[idx]
		s.localPeer.increment(favPeer.getRootHash(), favPeer.getPrevHash(), favPeer.getEntries())
		//log.Debug("New favourite block", "rootHash", string(favPeer.getRootHash()), "prevHash", string(favPeer.getPrevHash()))
	}

	// Only 1 favourite, have seen atlest 90% of expected hosts and they have all voted
	// for the single favourite
	if len(favourites) == 1 && len(s.peerMap) >= s.experimentParticipants && numVotes >= s.experimentParticipants && !s.converged {
		new := s.round() - s.prevRoundNumber
		s.convergeMap[string(favourites[0].getRootHash())] = new
		s.converged = true
		log.Info("Converged", "rounds", new)
	}

}

func (s *state) mergePeer(p *blockchain.PeerState) {
	if p.GetEntryHashes() == nil {
		return
	}

	if existingPeer := s.peerMap[p.GetId()]; existingPeer != nil {
		if p.GetEpoch() > existingPeer.getEpoch() {
			existingPeer.update(p.GetEpoch(), p.GetEntryHashes(), p.GetRootHash(), p.GetPrevHash())

			s.mergeMissingEntries(p.GetEntryHashes())
		}
	} else {
		newPeer := &peer{
			id:       p.GetId(),
			entries:  p.GetEntryHashes(),
			rootHash: p.GetRootHash(),
			prevHash: p.GetPrevHash(),
			epoch:    p.GetEpoch(),
			httpAddr: p.GetHttpAddr(),
		}

		s.peerMap[newPeer.id] = newPeer

		s.mergeMissingEntries(p.GetEntryHashes())
	}
}

func (s *state) mergeMissingEntries(entries [][]byte) {
	for _, e := range entries {
		key := string(e)
		if !s.pool.isPending(key) && !s.pool.isMissing(key) {
			s.pool.addMissing(key)
			log.Debug("Added to missing entries", "amount", len(s.pool.getAllMissing()))
		}
	}
}

func (s *state) add(e *entry) error {
	s.Lock()
	defer s.Unlock()

	return s.addNonLock(e)
}

func (s *state) addNonLock(e *entry) error {
	s.pool.addPending(e)

	err := s.inProgress.add(e)

	if err != nil && err != errFullBlock {
		log.Error(err.Error())
		return err
	} else if err == errFullBlock && !s.localPeer.hasFavourite() {
		s.localPeer.addBlock(s.inProgress)
	}

	return nil
}

func (s *state) mergeResponse(r *blockchain.StateResponse) {
	s.Lock()
	defer s.Unlock()

	if r == nil {
		log.Debug("Response is nil")
		return
	}

	for _, e := range r.GetEntries() {
		key := string(e.GetHash())
		if s.pool.isMissing(key) || s.pool.isPending(key) {
			log.Debug("Got missing entry in response", "amount", len(s.pool.getAllMissing()))
			s.pool.removeMissing(key)
			if !s.pool.isPending(key) {
				newEntry := &entry{data: e.GetContent(), hash: e.GetHash()}
				err := s.addNonLock(newEntry)
				if err != nil {
					log.Error(err.Error())
				}
			}
		}
	}
}

func (s *state) newRound() *block {
	s.Lock()
	defer s.Unlock()

	log.Debug("Picking next block")

	new := createBlock(s.localPeer.getPrevHash())

	if s.localPeer.getEntries() == nil {
		log.Error("no entries, adding empty block")
	}

	for _, e := range s.localPeer.getEntries() {
		key := string(e)
		ent := s.pool.getPending(key)
		if ent == nil {
			log.Error("Missing entry in fav block")
			continue
		}

		err := new.add(ent)
		if err != nil {
			log.Error(err.Error())
		} else {
			s.pool.removePending(key)
			s.pool.addConfirmed(ent)
		}
	}

	// TODO need a fallback when we we've not recieved all entries in the chosen block
	if eq := new.cmpRootHash(s.localPeer.getRootHash()); !eq {
		log.Error("Have different root hash for new block")
	}

	// TODO is this even correct?
	for _, p := range s.peerMap {
		p.reset()
	}

	s.pool.resetMissing()

	s.inProgress = createBlock(new.getRootHash())

	if full := s.pool.fillWithPending(s.inProgress); full {
		s.localPeer.addBlock(s.inProgress)
	}

	s.converged = false
	s.prevRoundNumber = s.round()

	return new
}

func (s *state) getHttpAddrs() []string {
	s.RLock()
	defer s.RUnlock()

	idx := 0
	ret := make([]string, len(s.peerMap))

	for _, p := range s.peerMap {
		ret[idx] = p.httpAddr
		idx++
	}

	return ret
}
