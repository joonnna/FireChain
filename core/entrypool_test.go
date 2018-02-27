package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PoolTestSuite struct {
	suite.Suite
	p *entryPool
}

func TestPoolTestSuite(t *testing.T) {
	suite.Run(t, new(PoolTestSuite))
}

func (suite *PoolTestSuite) SetupTest() {
	suite.p = newEntryPool()
}

func (suite *PoolTestSuite) TestAddPending() {
	e := &entry{hash: []byte("testhash")}
	k := string(e.hash)

	suite.p.addPending(e)

	assert.True(suite.T(), suite.p.isPending(k), "Valid add is not added to pending map")
	assert.True(suite.T(), suite.p.exists(k), "Valid add is not added to pending map")
	assert.NotEqual(suite.T(), nil, suite.p.getPending(k), "Failed to get pending entry after add")
}

func (suite *PoolTestSuite) TestAddConfirmed() {
	e := &entry{hash: []byte("testhash")}
	k := string(e.hash)

	suite.p.addConfirmed(e)

	assert.True(suite.T(), suite.p.isConfirmed(k), "Valid add is not added to pending map")
	assert.True(suite.T(), suite.p.exists(k), "Valid add is not added to pending map")
	assert.NotEqual(suite.T(), nil, suite.p.getConfirmed(k), "Failed to get confirmed entry after add")
}

func (suite *PoolTestSuite) TestRemovePending() {
	e := &entry{hash: []byte("testhash")}
	k := string(e.hash)

	suite.p.addPending(e)

	assert.True(suite.T(), suite.p.isPending(k), "Valid add is not added to pending map")

	suite.p.removePending(k)

	assert.False(suite.T(), suite.p.isPending(k), "Pending entry still present after removal")
	assert.False(suite.T(), suite.p.exists(k), "Pending entry still present after removal")
	assert.Nil(suite.T(), suite.p.getPending(k), "Pending entry still present after removal")
}

func (suite *PoolTestSuite) TestRemoveConfirmed() {
	e := &entry{hash: []byte("testhash")}
	k := string(e.hash)

	suite.p.addConfirmed(e)

	assert.True(suite.T(), suite.p.isConfirmed(k), "Valid add is not added to confrimed map")

	suite.p.removeConfirmed(k)

	assert.False(suite.T(), suite.p.isConfirmed(k), "Confirmed entry still present after removal")
	assert.False(suite.T(), suite.p.exists(k), "Confirmed entry still present after removal")
	assert.Nil(suite.T(), suite.p.getConfirmed(k), "Confirmed entry still present after removal")
}

func (suite *PoolTestSuite) TestFillWithPending() {
	e := &entry{data: []byte("testdata"), hash: []byte("testhash")}

	suite.p.addPending(e)

	b := createBlock(nil)

	assert.False(suite.T(), suite.p.fillWithPending(b), "Single entry should not fill block")

	for _, h := range b.getHashes() {
		assert.Equal(suite.T(), e.hash, h, "Entries dont match")
	}

	e2 := &entry{data: make([]byte, maxBlockSize+1), hash: []byte("testhash")}

	suite.p.addPending(e2)

	assert.True(suite.T(), suite.p.fillWithPending(b), "Entry with blocksize should fill block")
}

func (suite *PoolTestSuite) TestEmptyDiff() {
	for i := 0; i < 10; i++ {
		e := &entry{
			data: []byte(fmt.Sprintf("testdata%d", i)),
			hash: []byte(fmt.Sprintf("testhash%d", i)),
		}

		suite.p.addPending(e)
	}

	diff := suite.p.diff(nil)
	assert.Equal(suite.T(), len(suite.p.pending), len(diff), "Passing nil should return all pending")

	for _, e := range diff {
		assert.True(suite.T(), suite.p.isPending(string(e.GetHash())), "Content not equal")
	}
}

func (suite *PoolTestSuite) TestDiff() {
	other := make(map[string]bool)
	notShared := make(map[string]bool)

	numNotShared := 10

	for i := 0; i < numNotShared; i++ {
		e := &entry{
			data: []byte(fmt.Sprintf("testdata%d", i)),
			hash: []byte(fmt.Sprintf("testhash%d", i)),
		}
		notShared[string(e.hash)] = true
		suite.p.addPending(e)
	}

	for i := 15; i < 20; i++ {
		h := fmt.Sprintf("testhash%d", i)
		other[h] = true
	}

	diff := suite.p.diff(other)
	assert.Equal(suite.T(), len(diff), numNotShared, "Diff returns wrong amount of entries")

	for _, e := range diff {
		assert.NotNil(suite.T(), notShared[string(e.GetHash())], "Diff returns wrong entries")
	}

}
