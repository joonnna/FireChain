package core

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/joonnna/blocks/protobuf"
)

type ChainTestSuite struct {
	suite.Suite
	c *Chain
}

func TestChainTestSuite(t *testing.T) {
	suite.Run(t, new(ChainTestSuite))
}

func (suite *ChainTestSuite) SetupTest() {
	var err error
	suite.c, err = NewChain(nil)
	assert.NoError(suite.T(), err, "Failed to create chain")
}

/*
func (suite *ChainTestSuite) TearDownTest() {
}
*/
func (suite *ChainTestSuite) TestAdd() {
	err := suite.c.Add(nil)
	assert.Error(suite.T(), err, "Adding nil data return no error")

	err = suite.c.Add(make([]byte, 0))
	assert.Error(suite.T(), err, "Adding zero length data returns no error")

	err = suite.c.addBlock(nil)
	assert.Error(suite.T(), err, "Adding nil block returns no error")

	err = suite.c.addBlock(&block{})
	assert.NoError(suite.T(), err, "Adding non-nil block returns an error")

	assert.NotEqual(suite.T(), 0, suite.c.currBlock, "Current block is not updated after adding one")
}

func (suite *ChainTestSuite) TestHandleMsg() {
	_, err := suite.c.handleMsg(nil)
	assert.Error(suite.T(), err, "Nil message returns no error")

	m := &blockchain.State{
		PendingEntries: make(map[string]bool),
	}

	// Need to fill message with something before marshal
	m.PendingEntries["test"] = true

	bytes, err := proto.Marshal(m)
	assert.NoError(suite.T(), err, "Failed to marshal message")

	_, err = suite.c.handleMsg(bytes)
	assert.NoError(suite.T(), err, "Valid message returns error")

	_, err = suite.c.handleMsg([]byte("test"))
	assert.Error(suite.T(), err, "Message with random bytes returns no error")
}
