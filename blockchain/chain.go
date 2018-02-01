package blockchain

import (
	"errors"
	"sync"
)

var (
	errNoBlock = errors.New("No current block, panic?!")
)

type chain struct {
	blocks    map[uint64]*block
	currBlock uint64

	data      []int32
	dataMutex sync.RWMutex
}

func newChain() *chain {
	return &chain{
		blocks: make(map[uint64]*block),
	}
}

func (c *chain) Add(data int32) error {
	/*
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
	*/

	c.dataMutex.Lock()
	defer c.dataMutex.Unlock()

	exists := false

	for _, v := range c.data {
		if v == data {
			exists = true
			break
		}
	}

	if !exists {
		c.data = append(c.data, data)
	}

	return nil
}

func (c *chain) State() []int32 {
	c.dataMutex.RLock()
	defer c.dataMutex.RUnlock()

	buf := make([]int32, len(c.data))

	copy(buf, c.data)

	return buf
}
