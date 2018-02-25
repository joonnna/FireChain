package core

import (
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
