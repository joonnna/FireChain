package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Nil(suite.T(), diff, "Passing nil to diff should return nil")

	var empty [][]byte

	newDiff := suite.p.diff(empty)
	assert.Nil(suite.T(), newDiff, "Passing empty slice to diff should return nil")
}

func (suite *PoolTestSuite) TestDiff() {
	var other [][]byte

	notShared := make(map[string]bool)

	requestSize := 5
	local := 10

	for i := 0; i < local; i++ {
		e := &entry{
			hash: []byte(fmt.Sprintf("testhash%d", i)),
		}
		suite.p.addPending(e)
		if i < requestSize {
			other = append(other, e.hash)
			notShared[string(e.hash)] = true
		}
	}

	diff := suite.p.diff(other)
	assert.Equal(suite.T(), len(diff), requestSize, "Diff returns wrong amount of entries")

	for _, e := range diff {
		assert.NotNil(suite.T(), notShared[string(e.GetHash())], "Diff returns wrong entries")
	}
}

func (suite *PoolTestSuite) TestMissingToPendingNonFull() {
	e := &entry{hash: []byte("testhash")}

	key := string(e.hash)

	suite.p.addMissing(key)

	suite.p.missingToPending(e)

	assert.True(suite.T(), suite.p.isPending(key), "Entry is not present in pending")
}

func (suite *PoolTestSuite) TestMissingToPendingFull() {
	suite.p.cap = 1

	e1 := &entry{hash: []byte("testhash1"), data: []byte("1")}
	k1 := string(e1.hash)

	err := suite.p.addPending(e1)
	require.NoError(suite.T(), err, "Pending pool should not be full")

	require.True(suite.T(), suite.p.full, "Pending pool is not full")

	e2 := &entry{hash: []byte("testhash2"), data: []byte("2")}
	k2 := string(e2.hash)

	suite.p.addMissing(k2)

	suite.p.missingToPending(e2)

	assert.True(suite.T(), suite.p.isPending(k2), "Entry is not present in pending")
	assert.False(suite.T(), suite.p.isPending(k1), "Entry should be removed from pool")
	require.True(suite.T(), suite.p.full, "Pending pool is not full after entry switch")
}

func (suite *PoolTestSuite) TestMissingToPendingFavorite() {
	suite.p.cap = 1

	e1 := &entry{hash: []byte("testhash1"), data: []byte("1")}
	k1 := string(e1.hash)

	err := suite.p.addPending(e1)
	require.NoError(suite.T(), err, "Pending pool should not be full")

	suite.p.addFavorite(k1)

	require.True(suite.T(), suite.p.full, "Pending pool is not full")

	e2 := &entry{hash: []byte("testhash2"), data: []byte("2")}
	k2 := string(e2.hash)

	suite.p.addMissing(k2)

	suite.p.missingToPending(e2)

	assert.False(suite.T(), suite.p.isPending(k2), "Non-favorite entry switched replaced favorite entry")
	assert.True(suite.T(), suite.p.isPending(k1), "Favorite entry removed")
	require.True(suite.T(), suite.p.full, "Pending pool is not full after attempted switch")
}
