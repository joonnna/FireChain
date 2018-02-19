package core

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/cbergoon/merkletree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SearchTestSuite struct {
	suite.Suite
}

func TestSearchTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

func (suite *SearchTestSuite) TestEmptySlices() {
	var entries []*entry
	var tree []merkletree.Content

	e := &entry{}

	entries, tree = insertContent(entries, tree, e)
	assert.Equal(suite.T(), 1, len(entries), "Does not insert into entries when slice is empty")
	assert.Equal(suite.T(), 1, len(tree), "Does not insert into tree when slice is empty")
}

func (suite *SearchTestSuite) TestSorting() {
	var entries []*entry
	var tree []merkletree.Content

	var entries2 []*entry
	var tree2 []merkletree.Content

	e1 := &entry{data: []byte("random string....")}
	e2 := &entry{data: []byte("more random string yay")}
	e3 := &entry{data: []byte("best test case ever")}

	entries, tree = insertContent(entries, tree, e1)
	entries, tree = insertContent(entries, tree, e2)
	entries, tree = insertContent(entries, tree, e3)

	entries2, tree2 = insertContent(entries2, tree2, e3)
	entries2, tree2 = insertContent(entries2, tree2, e2)
	entries2, tree2 = insertContent(entries2, tree2, e1)

	fullTree, err := merkletree.NewTree(tree)
	assert.NoError(suite.T(), err, "Fail to build tree")

	fullTree2, err := merkletree.NewTree(tree2)
	assert.NoError(suite.T(), err, "Fail to build tree")

	fmt.Println(fullTree.MerkleRoot())
	fmt.Println(fullTree2.MerkleRoot())

	cmp := bytes.Equal(fullTree.MerkleRoot(), fullTree2.MerkleRoot())

	assert.True(suite.T(), cmp, "Merkle roots are not equal")
}
