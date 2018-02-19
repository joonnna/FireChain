package core

import (
	"bytes"
	"errors"

	"github.com/cbergoon/merkletree"
	"github.com/joonnna/blocks/protobuf"
)

var (
	errFullBlock         = errors.New("Block is full")
	errNoContent         = errors.New("Block blockchain has no content")
	errDifferentRootHash = errors.New("Block does not have same root hash")
	errAlreadyExists     = errors.New("Already exists in block")
)

// Written once and only read after, no need for mutexes

type block struct {
	//TODO change this
	treeContent []merkletree.Content
	entries     []*entry

	currSize uint32
	maxSize  uint32
	tree     *merkletree.MerkleTree

	rootHash []byte
	prevHash []byte
}

func createBlock(prevHash []byte) *block {
	return &block{
		maxSize:  128,
		prevHash: prevHash,
	}
}

func formBlock(data []byte) *block {
	return &block{}
}

func (b *block) addToBlock(e *entry) error {
	var err error

	size := uint32(len(e.data))

	if (b.currSize + size) >= b.maxSize {
		return errFullBlock
	}

	if exists := b.existsInBlock(e); exists {
		return errAlreadyExists
	}

	b.entries, b.treeContent = insertContent(b.entries, b.treeContent, e)

	if b.tree == nil {
		b.tree, err = merkletree.NewTree(b.treeContent)
		if err != nil {
			return err
		}
	} else {
		err = b.tree.RebuildTreeWith(b.treeContent)
		if err != nil {
			return err
		}
	}

	b.rootHash = b.tree.MerkleRoot()

	b.currSize += size

	return nil
}

func (b block) existsInBlock(new *entry) bool {
	for _, e := range b.entries {
		if eq := e.equal(new); eq {
			return true
		}
	}

	return false
}

func (b *block) add(e *entry) error {
	return b.addToBlock(e)
}

func (b *block) cmpRootHash(other []byte) bool {
	return bytes.Equal(b.rootHash, other)
}

func (b *block) cmpPrevHash(other []byte) bool {
	return bytes.Equal(b.prevHash, other)
}

func (b block) getPrevHash() []byte {
	return b.prevHash
}

// no need to copy, read only
func (b *block) getRootHash() []byte {
	return b.rootHash
	/*
		ret := make([]byte, len(b.rootHash))

		copy(ret, b.rootHash)

		return ret
	*/
}

func (b *block) content() [][]byte {
	ret := make([][]byte, len(b.entries))

	for idx, e := range b.entries {
		ret[idx] = e.data
	}

	return ret
}

func (b *block) getHashes() [][]byte {
	ret := make([][]byte, len(b.entries))

	for idx, e := range b.entries {
		ret[idx] = e.hash
	}

	return ret
}

func (b *block) blockToPbMsg() *blockchain.BlockContent {
	return &blockchain.BlockContent{
		RootHash: b.rootHash,
		PrevHash: b.prevHash,
		Content:  b.content(),
	}
}

/*
func pbMsgToBlock(m *blockchain.BlockContent) (*block, error) {
	content := m.GetContent()
	if content == nil {
		return nil, errNoContent
	}

	b := &block{
		prevHash: m.GetPrevHash(),
	}

	for _, data := range content {
		b.add(data)
	}

	if equal := b.cmpRootHash(m.GetRootHash()); !equal {
		return nil, errDifferentRootHash
	}
	return b, nil
}


func (b *block) mergeBlock(other [][]byte) error {
	b.Lock()
	defer b.Unlock()

	m := make(map[string]bool)

	for _, e := range b.entries {
		m[string(e.data)] = true
	}

	for _, d := range other {
		if _, exists := m[string(d)]; !exists {
			err := b.addToBlock(d)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
*/
