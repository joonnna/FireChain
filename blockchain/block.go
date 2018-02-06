package blockchain

import (
	"errors"
	"time"

	"github.com/joonnna/blocks/merkle"
)

var (
	errInvalidBlock = errors.New("Invalid block")
	errFullBlock    = errors.New("Block is full")
	maxSize         = 4096
)

type block struct {
	blockNum uint64
	currSize uint32

	prevHash []byte

	tree *merkle.Tree

	timestamp time.Time

	//entryMutex sync.RWMutex
	//entries    []*entry
}

func createBlock() *block {
	return &block{
		tree:      merkle.NewTree(),
		timestamp: time.Now(),
	}
}

func formBlock(data []byte) *block {
	return &block{}
}

func (b *block) add(data []byte) error {
	b.entryMutex.Lock()
	defer b.entryMutex.Unlock()

	size := len(data)

	if (b.currSize + size) >= maxSize {
		return errFullBlock
	}

	err := b.tree.Add(data)
	if err != nil {
		return err
	}

	b.currSize += size

	return nil
}
