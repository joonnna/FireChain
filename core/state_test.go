package core

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	log "github.com/inconshreveable/log15"
	"github.com/joonnna/blocks/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StateTestSuite struct {
	suite.Suite
	s *state
}

func TestStateTestSuite(t *testing.T) {
	r := log.Root()

	r.SetHandler(log.CallerFileHandler(log.StreamHandler(os.Stdout, log.TerminalFormat())))

	suite.Run(t, new(StateTestSuite))
}

func (suite *StateTestSuite) SetupTest() {
	suite.s = newState("local", "")
}

func (suite *StateTestSuite) TestEmptyMerge() {
	e := &entry{hash: []byte("testhash")}

	suite.s.pool.addPending(e)

	// Protobuf somehow fails without both fields populated
	m := &blockchain.State{
		Peers:          make(map[string]*blockchain.PeerState),
		PendingEntries: make(map[string]bool),
	}

	p := &peer{
		id:      "testpeer",
		entries: make([][]byte, 1),
	}

	suite.s.peerMap[p.id] = p

	resp, err := suite.s.merge(m)
	assert.NoError(suite.T(), err, "Merge fails")

	r := &blockchain.StateResponse{}

	err = proto.Unmarshal(resp, r)
	assert.NoError(suite.T(), err, "Unmarshal fails")

	entries := r.GetMissingEntries()
	assert.Equal(suite.T(), 1, len(entries), "Merge response does not contain entries")
	assert.Equal(suite.T(), e.hash, entries[0].GetHash(), "Merged entry has wrong hash")
}

func (suite *StateTestSuite) TestMerge() {
	// Protobuf somehow fails without both fields populated
	m := &blockchain.State{
		Peers:          make(map[string]*blockchain.PeerState),
		PendingEntries: make(map[string]bool),
	}

	p1 := &blockchain.PeerState{Id: "peer1", EntryHashes: make([][]byte, 1)}

	m.Peers[p1.GetId()] = p1

	_, err := suite.s.merge(m)
	assert.NoError(suite.T(), err, "Merge fails")

	// Local + merged peer, should be 2 entries in map
	assert.Equal(suite.T(), 2, len(suite.s.peerMap), "Peer not merged")
	assert.NotNil(suite.T(), suite.s.peerMap[p1.GetId()], "Peer not added to map")
}

func (suite *StateTestSuite) TestBlockVoting() {
	local := suite.s.peerMap["local"]

	local.rootHash = []byte("localfavblock")

	peers := make(map[string]*blockchain.PeerState)

	p1 := &blockchain.PeerState{
		Id:          "peer1",
		RootHash:    []byte("peer1favblock"),
		EntryHashes: make([][]byte, 1),
	}

	p2 := &blockchain.PeerState{
		Id:          "peer2",
		RootHash:    []byte("peer1favblock"),
		EntryHashes: make([][]byte, 1),
	}

	peers[p1.GetId()] = p1
	peers[p2.GetId()] = p2

	suite.s.mergePeers(peers)

	assert.Equal(suite.T(), p1.GetRootHash(), local.getRootHash(), "Favoured block is not chosen")
}

func (suite *StateTestSuite) TestMultipleFavourites() {
	local := suite.s.peerMap["local"]

	local.rootHash = []byte("localfavblock")

	peers := make(map[string]*blockchain.PeerState)

	p1 := &blockchain.PeerState{
		Id:          "peer1",
		RootHash:    []byte("peer1favblock"),
		EntryHashes: make([][]byte, 1),
	}

	p2 := &blockchain.PeerState{
		Id:          "peer2",
		RootHash:    []byte("peer1favblock"),
		EntryHashes: make([][]byte, 1),
	}

	p3 := &blockchain.PeerState{
		Id:          "peer3",
		RootHash:    []byte("localfavblock"),
		EntryHashes: make([][]byte, 1),
	}

	peers[p1.GetId()] = p1
	peers[p2.GetId()] = p2
	peers[p3.GetId()] = p3

	suite.s.mergePeers(peers)

	if eq := bytes.Equal(local.getRootHash(), p3.GetRootHash()); eq {
		assert.Equal(suite.T(), []byte("localfavblock"), local.getRootHash(), "Favoured block is not chosen")
	} else {
		assert.Equal(suite.T(), p1.GetRootHash(), local.getRootHash(), "Favoured block is not chosen")
	}

}

func (suite *StateTestSuite) TestMergeEmptyPeer() {
	local := suite.s.peerMap["local"]

	prevSize := len(suite.s.peerMap)

	local.update(uint64(1), make([][]byte, 1), make([]byte, 1), make([]byte, 1))

	prevEpoch := local.getEpoch()
	prevRoot := local.getRootHash()
	prevHash := local.getPrevHash()
	prevEntries := local.getEntries()

	p := &blockchain.PeerState{}

	suite.s.mergePeer(p)

	// Empty peer should not change any state
	assert.Equal(suite.T(), prevSize, len(suite.s.peerMap), "Empty peer added to map")
	assert.Equal(suite.T(), prevEpoch, local.getEpoch(), "Epoch changed")
	assert.Equal(suite.T(), prevRoot, local.getRootHash(), "Root hash changed")
	assert.Equal(suite.T(), prevHash, local.getPrevHash(), "Prev hash changed")
	assert.Equal(suite.T(), prevEntries, local.getEntries(), "Entries changed")
}

func (suite *StateTestSuite) TestMergePeer() {
	p := &blockchain.PeerState{
		Id:          "testpeer",
		EntryHashes: make([][]byte, 1),
		RootHash:    make([]byte, 2),
		PrevHash:    make([]byte, 3),
		Epoch:       uint64(1),
	}

	suite.s.mergePeer(p)

	newPeer := suite.s.peerMap[p.GetId()]

	// Local + merged peer, should be 2 entries in map
	assert.Equal(suite.T(), 2, len(suite.s.peerMap), "Merged peer not added to map")
	assert.NotNil(suite.T(), newPeer, "Merged peer is nil")

	assert.Equal(suite.T(), newPeer.id, p.GetId(), "Peer state not saved")
	assert.Equal(suite.T(), newPeer.getEntries(), p.GetEntryHashes(), "Peer state not saved")
	assert.Equal(suite.T(), newPeer.getRootHash(), p.GetRootHash(), "Peer state not saved")
	assert.Equal(suite.T(), newPeer.getPrevHash(), p.GetPrevHash(), "Peer state not saved")

	newEntryHashes := make([][]byte, 5)
	newRootHash := make([]byte, 10)
	newPrevHash := make([]byte, 15)
	newEpoch := p.GetEpoch() + 1

	p.EntryHashes = newEntryHashes
	p.RootHash = newRootHash
	p.PrevHash = newPrevHash
	p.Epoch = newEpoch

	suite.s.mergePeer(p)

	updatedPeer := suite.s.peerMap[p.GetId()]

	assert.Equal(suite.T(), newEntryHashes, updatedPeer.getEntries(), "Peer state not updated")
	assert.Equal(suite.T(), newRootHash, updatedPeer.getRootHash(), "Peer state not updated")
	assert.Equal(suite.T(), newPrevHash, updatedPeer.getPrevHash(), "Peer state not updated")

	newEntryHashes = make([][]byte, 25)
	newRootHash = make([]byte, 30)
	newPrevHash = make([]byte, 35)
	newEpoch = p.GetEpoch()

	p.EntryHashes = newEntryHashes
	p.RootHash = newRootHash
	p.PrevHash = newPrevHash
	p.Epoch = newEpoch

	suite.s.mergePeer(p)

	updatedPeer = suite.s.peerMap[p.GetId()]

	assert.NotEqual(suite.T(), newEntryHashes, updatedPeer.getEntries(), "Peer state updated with equal epoch")
	assert.NotEqual(suite.T(), newRootHash, updatedPeer.getRootHash(), "Peer state updated with equal epoch")
	assert.NotEqual(suite.T(), newPrevHash, updatedPeer.getPrevHash(), "Peer state updated with equal epoch")

}

func (suite *StateTestSuite) TestAdd() {
	e := &entry{data: []byte("testdata"), hash: []byte("testhash")}

	suite.s.add(e)

	local := suite.s.peerMap["local"]

	assert.False(suite.T(), local.hasFavourite(), "Should not have favourite with unfinished block")
	assert.Equal(suite.T(), 1, len(suite.s.inProgress.content()), "Entry not added to local progress block")

	e2 := &entry{data: make([]byte, maxBlockSize), hash: []byte("testhash2")}

	suite.s.add(e2)

	assert.True(suite.T(), local.hasFavourite(), "Should have favourite with finished block")
	assert.Equal(suite.T(), suite.s.localPeer.getRootHash(), suite.s.inProgress.getRootHash(), "Not same block stored in local peer and progress block")
}

func (suite *StateTestSuite) TestNewRound() {
	entries := make([]*entry, maxBlockSize)

	for i := 0; i < maxBlockSize; i++ {
		e := &entry{
			data: []byte(fmt.Sprintf("thisIsTestDataContentThatNeedsToBeALittleBig%d", i)),
			hash: []byte(fmt.Sprintf("testhash%d", i)),
		}
		suite.s.pool.addPending(e)
		entries[i] = e
	}

	local := suite.s.peerMap["local"]
	prevHash := local.getPrevHash()

	newBlock := suite.s.newRound()

	for _, h := range newBlock.getHashes() {
		assert.True(suite.T(), suite.s.pool.isConfirmed(string(h)), "Entry in block not confirmed")
		assert.False(suite.T(), suite.s.pool.isPending(string(h)), "Confirmed entry still pending")
	}

	for _, h := range suite.s.inProgress.getHashes() {
		assert.False(suite.T(), suite.s.pool.isConfirmed(string(h)), "Filled local block added to confirmed")
		assert.True(suite.T(), suite.s.pool.isPending(string(h)), "Filled local block removed from pending")
	}

	for _, h := range suite.s.localPeer.getEntries() {
		assert.False(suite.T(), suite.s.pool.isConfirmed(string(h)), "Filled local block added to confirmed")
		assert.True(suite.T(), suite.s.pool.isPending(string(h)), "Filled local block removed from pending")
	}

	assert.Equal(suite.T(), prevHash, newBlock.getRootHash(), "New block does not have correct prev")
	assert.Equal(suite.T(), suite.s.inProgress.getPrevHash(), newBlock.getRootHash(), "Progress block does not have new block as prev")
	assert.Equal(suite.T(), suite.s.localPeer.getPrevHash(), newBlock.getRootHash(), "Progress block does not have new block as prev")
}

func (suite *StateTestSuite) TestRoundReset() {
	for i := 0; i < 10; i++ {
		p := &peer{
			id:       fmt.Sprintf("test%d", i),
			rootHash: []byte(fmt.Sprintf("root%d", i)),
			prevHash: []byte(fmt.Sprintf("prev%d", i)),
			entries:  make([][]byte, 1),
		}
		suite.s.peerMap[p.id] = p
	}

	suite.s.newRound()

	for _, p := range suite.s.peerMap {
		assert.Nil(suite.T(), p.getRootHash(), "peer not reseted")
		assert.Nil(suite.T(), p.getPrevHash(), "peer not reseted")
		assert.Nil(suite.T(), p.getEntries(), "peer not reseted")
	}
}
