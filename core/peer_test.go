package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PeerTestSuite struct {
	suite.Suite
	p *peer
}

func TestPeerTestSuite(t *testing.T) {
	suite.Run(t, new(PeerTestSuite))
}

func (suite *PeerTestSuite) SetupTest() {
	suite.p = &peer{}
}

func (suite *PeerTestSuite) TestUpdate() {
	var epoch uint64 = 1

	entries := make([][]byte, 1)
	root := make([]byte, 2)
	prev := make([]byte, 3)

	suite.p.update(epoch, entries, root, prev)

	assert.Equal(suite.T(), suite.p.getEpoch(), epoch, "Epoch not updated")

	// These assertions will compare slice content/len, not pointer values.
	assert.Equal(suite.T(), suite.p.getEntries(), entries, "Entries not updated")
	assert.Equal(suite.T(), suite.p.getRootHash(), root, "Roothash not updated")
	assert.Equal(suite.T(), suite.p.getPrevHash(), prev, "PrevHash not updated")
}

func (suite *PeerTestSuite) TestAddBlock() {
	data := make([]byte, 10)
	prev := make([]byte, 1)

	b := createBlock(prev)

	err := b.add(&entry{data: data})
	assert.NoError(suite.T(), err, "Failed to add to block")

	suite.p.addBlock(b)
	assert.Equal(suite.T(), uint64(1), suite.p.getEpoch(), "Epoch not updated")

	// These assertions will compare slice content/len, not pointer values.
	assert.Equal(suite.T(), b.getHashes(), suite.p.getEntries(), "Entries not updated")
	assert.Equal(suite.T(), b.getRootHash(), suite.p.getRootHash(), "Roothash not updated")
	assert.Equal(suite.T(), b.getPrevHash(), suite.p.getPrevHash(), "PrevHash not updated")
}

func (suite *PeerTestSuite) TestHasFavourite() {
	data := make([]byte, 10)
	prev := make([]byte, 1)

	b := createBlock(prev)

	err := b.add(&entry{data: data})
	assert.NoError(suite.T(), err, "Failed to add to block")

	suite.p.addBlock(b)

	assert.True(suite.T(), suite.p.hasFavourite(), "Peer does not have favourite after block add")
}

func (suite *PeerTestSuite) TestIncrement() {
	entries := make([][]byte, 1)
	root := make([]byte, 2)
	prev := make([]byte, 3)

	suite.p.increment(root, prev, entries)

	assert.Equal(suite.T(), uint64(1), suite.p.getEpoch(), "Epoch not updated")

	// These assertions will compare slice content/len, not pointer values.
	assert.Equal(suite.T(), suite.p.getEntries(), entries, "Entries not updated")
	assert.Equal(suite.T(), suite.p.getRootHash(), root, "Roothash not updated")
	assert.Equal(suite.T(), suite.p.getPrevHash(), prev, "PrevHash not updated")
}

func (suite *PeerTestSuite) TestReset() {
	data := make([]byte, 10)
	prev := make([]byte, 1)

	b := createBlock(prev)

	err := b.add(&entry{data: data})
	assert.NoError(suite.T(), err, "Failed to add to block")

	suite.p.addBlock(b)

	suite.p.reset()

	assert.Nil(suite.T(), suite.p.getRootHash(), "Failed to reset root hash")
	assert.Nil(suite.T(), suite.p.getPrevHash(), "Failed to reset prev hash")
	assert.Nil(suite.T(), suite.p.getEntries(), "Failed to reset entries hash")
	assert.Equal(suite.T(), uint64(1), suite.p.getEpoch(), "Epoch gets reset, should not")
}
