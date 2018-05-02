package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"sync"

	"github.com/golang/protobuf/proto"

	log "github.com/inconshreveable/log15"
	"github.com/joonnna/blocks/protobuf"
	"github.com/joonnna/ifrit"
)

var (
	errNoSignature = errors.New("Message contained nil signature")
	errNoHashes    = errors.New("Peer had no hashes")
	errNoPeerState = errors.New("No peer table in request")
)

type state struct {
	sync.RWMutex
	peerMap    map[string]*peer
	pool       *entryPool
	localPeer  *peer
	inProgress *block

	ifrit *ifrit.Client

	// DO NOT SET TO TRUE, only for skipping signature checks in test code
	skipSignatures bool

	// Experiments only
	prevRoundNumber        uint32
	convergeMap            map[string]uint32
	converged              bool
	experimentParticipants uint32
}

func newState(ifrit *ifrit.Client, httpAddr string, hosts uint32) *state {
	peerMap := make(map[string]*peer)
	localPeer := &peer{
		id:       ifrit.Id(),
		httpAddr: httpAddr,
	}
	peerMap[localPeer.id] = localPeer

	return &state{
		peerMap:    peerMap,
		pool:       newEntryPool(),
		localPeer:  localPeer,
		inProgress: createBlock(nil),
		ifrit:      ifrit,
		experimentParticipants: hosts,
		convergeMap:            make(map[string]uint32),
	}
}

func (s *state) bytes() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	m := &blockchain.State{
		Peers:          make(map[string]*blockchain.PeerInfo),
		MissingEntries: s.pool.getAllMissing(),
	}
	for _, p := range s.peerMap {
		s := &blockchain.PeerInfo{
			Id:       p.id,
			Epoch:    p.getEpoch(),
			RootHash: p.getRootHash(),
		}

		m.Peers[s.GetId()] = s
	}

	return proto.Marshal(m)
}

func (s *state) diff(other *blockchain.State) ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	resp := &blockchain.StateResponse{}

	if peers := other.GetPeers(); peers != nil {
		resp.Peers = s.peersDiff(peers)
	} else {
		log.Error("No peers state, returning error")
		return nil, errNoPeerState
	}

	if entries := other.GetMissingEntries(); entries != nil {
		resp.Entries = s.pool.diff(entries)
		log.Debug("Request for missing entries", "amount", len(entries))
	}

	b, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *state) peersDiff(peers map[string]*blockchain.PeerInfo) []*blockchain.PeerState {
	var diff []*blockchain.PeerState

	for _, p := range s.peerMap {
		if remotePeer, ok := peers[p.id]; ok {
			if p.getEpoch() > remotePeer.GetEpoch() && !bytes.Equal(p.getRootHash(), remotePeer.GetRootHash()) {
				diff = append(diff, p.toPbMsg())
			}
		} else {
			diff = append(diff, p.toPbMsg())
		}
	}

	return diff
}

func (s *state) mergePeers(peers []*blockchain.PeerState) {
	var favourites []*peer
	blockVotes := make(map[string]int)

	numVotes := 0

	for _, p := range peers {
		if p.GetId() == s.localPeer.id {
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

	log.Debug("Convergence progress", "ifrit_alive", len(s.ifrit.Members()), "peerMap", len(s.peerMap), "votes", numVotes)

	if numFavs := len(favourites); numFavs >= 1 {
		idx := 0
		if numFavs > 1 {
			idx = rand.Int() % numFavs
		}

		favPeer := favourites[idx]

		equal := s.localPeer.isEqual(favPeer)

		// favoured the same block already, pools are already correct
		if equal {
			return
		}

		s.localPeer.epoch++

		s.updateLocalPeer(favPeer.getRootHash(), favPeer.getPrevHash(), favPeer.getEntries())

		s.pool.resetFavorite()
		s.pool.resetMissing()
		for _, h := range favPeer.getEntries() {
			key := string(h)
			s.pool.addFavorite(key)
			if !s.pool.isPending(key) {
				s.pool.addMissing(key)
			}
		}
		k := string(favPeer.getRootHash())
		s.convergeMap[k] = (s.ifrit.GetGossipRounds() - s.prevRoundNumber)
	}
	//log.Debug("New favourite block", "rootHash", string(favPeer.getRootHash()), "prevHash", string(favPeer.getPrevHash()))
	//log.Debug("Convergence progress", "ifrit_alive", participants, "peerMap", len(s.peerMap), "votes", numVotes)
}

func (s *state) mergePeer(p *blockchain.PeerState) {
	var r, sign []byte

	if p == nil {
		log.Debug("nil peer")
		return
	}

	id := p.GetId()

	if p.GetEntryHashes() == nil {
		log.Debug(errNoHashes.Error())
		return
	}

	if existingPeer, ok := s.peerMap[id]; ok {
		if p.GetEpoch() > existingPeer.getEpoch() {
			if !s.skipSignatures {
				signature := p.GetSignature()
				if signature == nil {
					log.Error(errNoSignature.Error())
					return
				}

				r = signature.GetR()
				sign = signature.GetS()

				if valid := s.checkSignature(p, r, sign); !valid {
					return
				}
			}
			existingPeer.update(p.GetEpoch(), p.GetEntryHashes(), p.GetRootHash(), p.GetPrevHash(), r, sign)
		}
	} else {
		if !s.skipSignatures {
			signature := p.GetSignature()
			if signature == nil {
				log.Error(errNoSignature.Error())
				return
			}

			r = signature.GetR()
			sign = signature.GetS()

			if valid := s.checkSignature(p, r, sign); !valid {
				return
			}
		}
		newPeer := &peer{
			id:       id,
			entries:  p.GetEntryHashes(),
			rootHash: p.GetRootHash(),
			prevHash: p.GetPrevHash(),
			epoch:    p.GetEpoch(),
			httpAddr: p.GetHttpAddr(),
			r:        r,
			s:        sign,
		}

		s.peerMap[id] = newPeer
	}
}

/*
func (s *state) mergeMissingEntries(entries [][]byte) {
	for _, e := range entries {
		key := string(e)
		if !s.pool.isPending(key) && !s.pool.isMissing(key) {
			s.pool.addMissing(key)
			log.Debug("Added to missing entries", "amount", len(s.pool.getAllMissing()))
		}
	}
}
*/

func (s *state) add(e *entry) error {
	s.Lock()
	defer s.Unlock()

	return s.addNonLock(e)
}

func (s *state) addNonLock(e *entry) error {
	err := s.pool.addPending(e)
	if err != nil {
		return err
	}

	err = s.inProgress.add(e)
	if err != nil && err != errFullBlock {
		log.Error(err.Error())
		return err
	} else if err == errFullBlock && !s.localPeer.hasFavourite() {
		s.localPeer.addBlock(s.inProgress)
		err := s.sign(s.localPeer)
		if err != nil {
			log.Error(err.Error())
		}

		key := string(s.localPeer.getRootHash())
		s.convergeMap[key] = (s.ifrit.GetGossipRounds() - s.prevRoundNumber)
	}

	return nil
}

func (s *state) mergeResponse(r *blockchain.StateResponse) {
	s.Lock()
	defer s.Unlock()

	if r == nil {
		log.Error("Response is nil")
		return
	}

	for _, e := range r.GetEntries() {
		key := string(e.GetHash())
		if s.pool.isMissing(key) {
			if !s.pool.isPending(key) {
				newEntry := &entry{data: e.GetContent(), hash: e.GetHash()}
				s.pool.missingToPending(newEntry)

				log.Debug("Got missing entry in response, still missing", "amount", len(s.pool.getAllMissing()))
			}
		}
	}

	s.mergePeers(r.GetPeers())
}

func (s *state) newRound() (*block, []byte) {
	s.Lock()
	defer s.Unlock()

	log.Debug("Picking next block")

	convergeValues := s.getConvergeValue(string(s.localPeer.getRootHash()))

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

	s.prevRoundNumber = s.ifrit.GetGossipRounds()
	s.convergeMap = nil
	s.convergeMap = make(map[string]uint32)

	s.pool.fillWithPending(s.inProgress)
	s.localPeer.addBlock(s.inProgress)
	err := s.sign(s.localPeer)
	if err != nil {
		log.Error(err.Error())
	}

	key := string(s.localPeer.getRootHash())
	s.convergeMap[key] = 0

	return new, convergeValues
}

func (s *state) sign(p *peer) error {
	msg := &blockchain.PeerState{
		Id:          p.id,
		EntryHashes: p.getEntries(),
		RootHash:    p.getRootHash(),
		PrevHash:    p.getPrevHash(),
		Epoch:       p.getEpoch(),
		HttpAddr:    p.httpAddr,
	}

	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	r, sign, err := s.ifrit.Sign(bytes)
	if err != nil {
		return err
	}

	p.addSignature(r, sign)

	return nil
}

func (s *state) checkSignature(p *blockchain.PeerState, r, sign []byte) bool {
	if r == nil || sign == nil {
		log.Error(errNoSignature.Error())
		return false
	}

	p.Signature = nil

	bytes, err := proto.Marshal(p)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	return s.ifrit.VerifySignature(r, sign, bytes, p.GetId())
}

func (s *state) getHttpAddrs() []string {
	s.RLock()
	defer s.RUnlock()

	ret := make([]string, 0, len(s.peerMap))

	for _, p := range s.peerMap {
		ret = append(ret, p.httpAddr)
	}

	return ret
}

func (s *state) updateLocalPeer(rootHash, prevHash []byte, entries [][]byte) {
	s.localPeer.entries = entries
	s.localPeer.rootHash = rootHash
	s.localPeer.prevHash = prevHash

	err := s.sign(s.localPeer)
	if err != nil {
		log.Error(err.Error())
	}
}

func (s *state) getConvergeValue(hash string) []byte {
	conv, ok := s.convergeMap[hash]
	if !ok {
		log.Error("Convergence not found", "hash", hash)
		for k, v := range s.convergeMap {
			log.Debug("Map entries", "hash", k, "rounds", v)
		}

		return nil
	}

	var b bytes.Buffer

	res := struct {
		Hash     string
		Converge uint32
	}{
		hash,
		conv,
	}

	err := json.NewEncoder(&b).Encode(res)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return b.Bytes()
}
