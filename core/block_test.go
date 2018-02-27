package core

import (
	"bytes"
	"os"
	"testing"

	log "github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BlockTestSuite struct {
	suite.Suite
}

func TestBlockTestSuite(t *testing.T) {
	r := log.Root()

	r.SetHandler(log.CallerFileHandler(log.StreamHandler(os.Stdout, log.TerminalFormat())))

	suite.Run(t, new(BlockTestSuite))
}

func (suite *BlockTestSuite) TestAdd() {
	b := createBlock(nil)
	e1 := &entry{}

	err := b.addToBlock(e1)
	assert.Error(suite.T(), err, "Adding data with nil data does not return an error")

	e2 := &entry{data: []byte("test")}
	e3 := &entry{data: []byte("test")}

	err = b.addToBlock(e2)
	assert.NoError(suite.T(), err, "Valid add fails")

	r := b.getRootHash()
	assert.NotEqual(suite.T(), nil, r, "Root hash is nil after successful add")

	assert.NotEqual(suite.T(), 0, b.currSize, "Current size if 0 after successful add")

	assert.NotEqual(suite.T(), nil, b.entries, "Block entries are nil after successful add")
	assert.NotEqual(suite.T(), nil, b.treeContent, "Tree content are nil after successful add")

	content := b.content()
	assert.Equal(suite.T(), 1, len(content), "After 1 add, content should have length 1")
	assert.True(suite.T(), bytes.Equal(e2.data, content[0]), "Added content should be equal")

	err = b.addToBlock(e3)
	assert.Error(suite.T(), err, "Adding same entry twice returns non-nil error")

	bytes := make([]byte, maxBlockSize+1)

	e4 := &entry{data: bytes}
	err = b.addToBlock(e4)
	assert.Error(suite.T(), err, "Adding too big entry returns non-nil error")
}
