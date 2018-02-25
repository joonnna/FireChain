package core

import (
	"crypto/sha256"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/joonnna/blocks/protobuf"
	"github.com/joonnna/ifrit"

	log "github.com/inconshreveable/log15"
)

var (
	errNilBlock           = errors.New("Tried to add nil block")
	errNoBlock            = errors.New("No current block, panic?!")
	errInvalidBlockHeader = errors.New("Block header contained empty fields")
	errInvalidBlock       = errors.New("Block contains empty fields")
	errDifferentPrev      = errors.New("Different previous block")
	errNoLocalPeer        = errors.New("Can't find local peer representation")
	errNoMsgContent       = errors.New("Message content is either nil or of zero length")
	ErrNoData             = errors.New("Given data is of zero length")
)

type Chain struct {
	blockMapMutex sync.RWMutex
	blocks        map[uint64]*block
	currBlock     uint64

	queueMutex sync.RWMutex
	queuedData [][]byte

	state *state

	ifrit *ifrit.Client
	id    string

	httpListener net.Listener

	exitChan chan bool
}

func NewChain(conf *ifrit.Config) (*Chain, error) {
	i, err := ifrit.NewClient(conf)
	if err != nil {
		return nil, err
	}

	l, err := initHttp()
	if err != nil {
		return nil, err
	}

	return &Chain{
		blocks:       make(map[uint64]*block),
		state:        newState(i.Id(), l.Addr().String()),
		ifrit:        i,
		httpListener: l,
		exitChan:     make(chan bool),
	}, nil
}

func (c *Chain) Start() {
	c.ifrit.RegisterMsgHandler(c.handleMsg)
	c.ifrit.RegisterResponseHandler(c.handleResponse)
	go c.ifrit.Start()
	go c.startHttp()

	c.blockLoop()
}

// Returns the address of the underlying ifrit client.
func (c Chain) Addr() string {
	return c.ifrit.Addr()
}

// Returns the address of the http module of the chain, used for worm input.
func (c Chain) HttpAddr() string {
	return c.httpListener.Addr().String()
}

func (c *Chain) Stop() {
	c.ifrit.ShutDown()
	c.stopHttp()

	close(c.exitChan)
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
	if data == nil || len(data) <= 0 {
		return nil, errNoMsgContent
	}

	msg := &blockchain.State{}

	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp, err := c.state.merge(msg)
	if err != nil {
		log.Error(err.Error())
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
		log.Error(err.Error())
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
		select {
		case <-c.exitChan:
			return

		case <-time.After(time.Second * 600):
			c.pickFavouriteBlock()
			c.updateState()
		}
	}
}

func (c *Chain) pickFavouriteBlock() {
	newBlock := c.state.newRound()

	// TODO if we get error here, we need a fallback strat
	err := c.addBlock(newBlock)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("New block", "rootHash", string(newBlock.getRootHash()), "prevhash", string(newBlock.getPrevHash()))
}

func (c *Chain) updateState() {
	bytes, err := c.state.bytes()
	if err != nil {
		log.Error(err.Error())
		return
	}
	c.ifrit.SetGossipContent(bytes)
}

func (c *Chain) addBlock(b *block) error {
	c.blockMapMutex.Lock()
	defer c.blockMapMutex.Unlock()

	if b == nil {
		return errNilBlock
	}

	c.currBlock++
	c.blocks[c.currBlock] = b

	return nil
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
			log.Error(err.Error())
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
			log.Error(err.Error())
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
