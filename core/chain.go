package core

import (
	"crypto/sha256"
	"errors"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/joonnna/blocks/protobuf"
	"github.com/joonnna/ifrit"
	"github.com/joonnna/ifrit/logger"
)

var (
	errNoBlock            = errors.New("No current block, panic?!")
	errInvalidBlockHeader = errors.New("Block header contained empty fields")
	errInvalidBlock       = errors.New("Block contains empty fields")
	errDifferentPrev      = errors.New("Different previous block")
	errNoLocalPeer        = errors.New("Can't find local peer representation")
	ErrNoData             = errors.New("Given data is of zero length")
)

type Chain struct {
	blockMapMutex sync.RWMutex
	blocks        map[uint64]*block
	currBlock     uint64

	localBlockMutex sync.RWMutex
	localBlock      *block

	queueMutex sync.RWMutex
	queuedData [][]byte

	state *state

	ifrit *ifrit.Client
	id    string

	log *logger.Log
}

func NewChain(entryAddr string, log *logger.Log) (*Chain, error) {
	i, err := ifrit.NewClient(entryAddr, &ifrit.Config{Visualizer: true, VisAddr: "127.0.0.1:8095"})
	if err != nil {
		return nil, err
	}

	return &Chain{
		blocks: make(map[uint64]*block),
		state:  newState(i.Id(), log),
		ifrit:  i,
		log:    log,
	}, nil
}

func (c *Chain) Start() {
	c.ifrit.RegisterMsgHandler(c.handleMsg)
	c.ifrit.RegisterResponseHandler(c.handleResponse)
	go c.ifrit.Start()

	c.blockLoop()
}

func (c *Chain) Add(data []byte) error {
	if data == nil || len(data) == 0 {
		return ErrNoData
	}

	e := &entry{
		data: data,
		hash: hashBytes(data),
	}

	c.state.add(e)

	c.updateState()

	return nil
}

func (c *Chain) handleMsg(data []byte) ([]byte, error) {
	msg := &blockchain.State{}

	err := proto.Unmarshal(data, msg)
	if err != nil {
		c.log.Err.Println(err)
		return nil, err
	}

	resp, err := c.state.merge(msg)
	if err != nil {
		c.log.Err.Println(err)
		return nil, err
	}

	c.updateState()

	return resp, nil
}

func (c *Chain) handleResponse(resp []byte) {
	if resp == nil || len(resp) == 0 {
		return
	}
	msg := &blockchain.StateResponse{}

	err := proto.Unmarshal(resp, msg)
	if err != nil {
		c.log.Err.Println(err)
		return
	}

	//c.log.Debug.Println("Entries in response: ", len(msg.GetMissingEntries()))

	for _, e := range msg.GetMissingEntries() {
		key := string(e.GetHash())
		if exists := c.state.exists(key); !exists {
			newEntry := &entry{data: e.GetContent(), hash: e.GetHash()}
			c.state.add(newEntry)
		}
	}

	c.updateState()
}

func (c *Chain) blockLoop() {
	for {
		time.Sleep(time.Second * 60)
		c.pickFavouriteBlock()
		c.updateState()
	}
}

func (c *Chain) pickFavouriteBlock() {
	var prev []byte

	curr := c.getCurrBlock()

	if curr != nil {
		prev = curr.getPrevHash()
	}

	newBlock := c.state.newRound(prev)

	c.addBlock(newBlock)
	c.log.Info.Println("Favourite block is now: ", string(newBlock.getRootHash()))
	c.log.Info.Println("Prev is: ", string(newBlock.getPrevHash()))
}

func (c *Chain) updateState() {
	bytes, err := c.state.bytes()
	if err != nil {
		c.log.Err.Println(err)
		return
	}
	c.ifrit.SetGossipContent(bytes)
}

func (c *Chain) addBlock(b *block) {
	c.blockMapMutex.Lock()
	defer c.blockMapMutex.Unlock()

	c.currBlock++
	c.blocks[c.currBlock] = b
}

func (c *Chain) getCurrBlock() *block {
	c.blockMapMutex.RLock()
	defer c.blockMapMutex.RUnlock()

	return c.blocks[c.currBlock]
}

func (c *Chain) addToQueue(data []byte) {
	c.queueMutex.Lock()
	defer c.queueMutex.Unlock()

	c.queuedData = append(c.queuedData, data)
}

func hashBytes(data []byte) []byte {
	return sha256.New().Sum(data)
}

/*
func (c *Chain) generateResponse(msg *blockchain.BlockHeader) ([]byte, error) {
	var bytes []byte
	var err error

	if msg == nil {
		c.log.Info.Println("Got empty block header")
		return nil, nil
	}

	rootHash := msg.GetRootHash()
	prevHash := msg.GetPrevHash()

	if rootHash == nil || prevHash == nil {
		return nil, errInvalidBlockHeader
	}

	if b, ok := c.blocks[c.getCurrBlock()]; !ok {
		c.log.Err.Println(errNoBlock)
		return nil, errNoBlock
	} else if prevEqual := b.cmpPrevHash(prevHash); !prevEqual {
		return nil, errDifferentPrev
	} else if rootEqual := b.cmpRootHash(rootHash); !rootEqual {
		resp := &blockchain.StateResponse{
			Block: b.blockToPbMsg(),
		}
		bytes, err = proto.Marshal(resp)
		if err != nil {
			c.log.Err.Println(err)
			return nil, err
		}

		return bytes, nil
	} else {
		return nil, nil
	}
}



func (c *Chain) mergeCurrBlock(msg *blockchain.BlockContent) {
	rootHash := msg.GetRootHash()
	prevHash := msg.GetPrevHash()
	content := msg.GetContent()

	if rootHash == nil || prevHash == nil || content == nil {
		c.log.Err.Println(errInvalidBlock)
		return
	}

	if b, ok := c.blocks[c.getCurrBlock()]; !ok {
		c.log.Err.Println(errNoBlock)
	} else {
		if equal := b.cmpPrevHash(prevHash); !equal {
			c.log.Err.Println(errDifferentPrev)
			return
		}

		if equal := b.cmpRootHash(rootHash); equal {
			c.log.Info.Println("Received an equal block in response")
			return
		}

		err := b.mergeBlock(content)
		if err != nil {
			c.log.Err.Println(err)
		}
	}

}


func (c *Chain) createLocalBlock(prev []byte) *block {
	c.localBlockMutex.Lock()
	defer c.localBlockMutex.Unlock()

	c.localBlock = createBlock(prev)

	return c.localBlock
}

func (c *Chain) getLocalBlock() *block {
	c.localBlockMutex.RLock()
	defer c.localBlockMutex.RUnlock()

	return c.localBlock
}

func (c *Chain) fillLocalBlock() {
	c.localBlockMutex.RLock()
	defer c.localBlockMutex.RUnlock()

	c.pool.fillWithPending(c.localBlock)
}
*/
