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

	peerMap map[string]*peer

	//data      []byte
	//dataMutex sync.RWMutex
}

type peer struct {
	id        string
	favBlock  []byte
	epoch     uint64
	signature []byte
}

func newChain() *chain {
	return &chain{
		blocks: make(map[uint64]*block),
	}
}

func (c *chain) handleMsg(data []byte) ([]byte, error) {
	for p, id := range data {
		if existingPeer, exists := c.peerMap[id]; !exists {
			add()
		}

		if p.block != existingPeer.block {
			changeBlock()
		}
	}

	incrementEpoch()
}

func (c *chain) Add(data []byte) error {
	if b, ok := c.blocks[c.currBlock]; !ok {
		return errNoBlock
	}

	err := b.add(data)
	if err == errFullBlock {
		newBlock := createBlock()
		c.addBlock(newBlock)
		return newBlock.add(data)
	} else if err != nil {
		return err
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

func (c *chain) addBlock(b *block) {
	c.blockMutex.Lock()
	defer c.blockMutex.Unlock()

	c.currBlock++
	c.blocks[c.currBlock] = b
}
