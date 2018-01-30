package blockchain

import (
	"errors"
)

var (
	errNoBlock = errors.New("No current block, panic?!")
)

type chain struct {
	blocks    map[uint64]*block
	currBlock uint64
}

func (c *chain) newChain() *chain {
	return &chain{
		blocks: make(map[uint64]*block),
	}
}

func (c *chain) Add(data []byte) error {
	if b, ok := c.blocks[c.currBlock]; !ok {
		return errNoBlock
	}

	err := b.add(data)
	if err == errFullBlock {
		newBlock := createBlock()
		c.currBlock++
		c.blocks[c.currBlock] = newBlock
		return c.add(data)
	} else if err != nil {
		return err
	}

	return nil
}

func (c *chain) state() error {

}
