package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"

	"github.com/cbergoon/merkletree"
)

var (
	errInvalidBlock = errors.New("Invalid block")
	errFullBlock    = errors.New("Block is full")
)

type block struct {
	contentMutex sync.RWMutex
	entries      []*entry
	currSize     uint32
	dirty        bool
	tree         *merkletree.MerkleTree
	headerHash   []byte

	*header
}

type header struct {
	merkleHash []byte
	prevHash   []byte
	targetSize uint32
}

func createBlock() *block {
	return &block{}
}

func formBlock(data []byte) *block {
	return &block{}
}

func (b *block) add(data []byte) error {
	b.contentMutex.Lock()
	defer b.contentMutex.Unlock()

	size := len(data)

	if (b.currSize + size) >= maxSize {
		return errFullBlock
	}

	b.entries = append(b.entries, &entry{data: data})

	if b.tree == nil {
		b.tree, err = merkletree.NewTree(b.entries)
		if err != nil {
			return err
		}
	} else {
		err := b.tree.RebuildTreeWith(b.entries)
		if err != nil {
			return err
		}
	}

	b.dirty = true

	b.currSize += size

	return nil
}

func (b *block) headerHash() []byte {
	b.contentMutex.Lock()
	defer b.contentMutex.Unlock()

	var src []byte

	if b.dirty {
		src = hashHeader(b.header)
		b.headerHash = src
	} else {
		src = b.headerHash
	}

	ret := make([]byte, len(src))

	copy(ret, src)

	b.dirty = false

	return ret

}

func hashHeader(h *header) []byte {
	//TODO Change this
	return sha256.Sum256(fmt.Sprintf("%v", h))
}
