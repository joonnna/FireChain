package core

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/joonnna/blocks/protobuf"
	"github.com/joonnna/ifrit"
	"github.com/joonnna/ifrit/netutil"

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
	errNoChainState       = errors.New("No chain state in request")
)

type Chain struct {
	blockMapMutex sync.RWMutex
	blocks        map[uint64]*block
	currBlock     uint64
	blockPeriod   time.Duration

	queueMutex sync.RWMutex
	queuedData [][]byte

	state *state

	ifrit *ifrit.Client
	id    string

	httpListener net.Listener

	exitChan chan bool

	forkCheckTimeout time.Duration

	// Experiments only
	saturationTimeout time.Duration
	experiment        bool
	expAddr           string
	httpAddr          string
	ExpChan           chan bool
	expStarted        bool
	expMutex          sync.Mutex
}

/*
type ExpArgs struct {
	SaturationTimeout time.Duration
	NumHosts          uint32
}
*/
/*
type test interface {
	Bytes()
	Merge()
	MergeResponse()
	NewRound()
	Add()
}
*/

func (c *Chain) fillFirstBlock() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < int(blockSize/entrySize)+1; i++ {
		buf := make([]byte, entrySize)
		_, err := rand.Read(buf)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		err = c.Add(buf)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
}

func NewChain(conf *ifrit.Config, expAddr string, forkLimit float64, blockPeriod, bSize, eSize uint32) (*Chain, error) {
	entrySize = eSize
	blockSize = bSize * 1024

	l, err := initHttp()
	if err != nil {
		return nil, err
	}

	ip := netutil.GetLocalIP()
	port := strings.Split(l.Addr().String(), ":")[1]
	httpAddr := fmt.Sprintf("%s:%s", ip, port)

	conf.VisAppAddr = fmt.Sprintf("http://%s", httpAddr)
	i, err := ifrit.NewClient(conf)
	if err != nil {
		return nil, err
	}

	exp := false
	if expAddr != "" {
		exp = true
	}

	return &Chain{
		blocks:           make(map[uint64]*block),
		state:            newState(i, httpAddr, forkLimit),
		ifrit:            i,
		httpListener:     l,
		exitChan:         make(chan bool),
		blockPeriod:      time.Minute * time.Duration(blockPeriod),
		ExpChan:          make(chan bool),
		experiment:       exp,
		expAddr:          expAddr,
		httpAddr:         httpAddr,
		forkCheckTimeout: time.Minute * 1,
	}, nil
}

func (c *Chain) Start() {
	go c.ifrit.Start()
	go c.startHttp()

	<-c.ExpChan
	c.fillFirstBlock()
	c.ifrit.RecordGossipRounds()
	c.ifrit.RegisterGossipHandler(c.handleMsg)
	c.ifrit.RegisterResponseHandler(c.handleResponse)
	c.ifrit.RegisterResponseHandler(c.handleResponse)
	c.ifrit.RegisterMsgHandler(c.handleForkMsg)
	go c.forkResolver()
	c.blockLoop()
}

// Returns the address of the underlying ifrit client.
func (c Chain) Addr() string {
	return c.ifrit.Addr()
}

// Returns the address of the http module of the chain, used for worm input.
func (c Chain) HttpAddr() string {
	//return c.httpListener.Addr().String()
	return c.httpAddr
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

	err := c.state.add(e)
	if err != nil {
		return err
	}

	c.updateState()

	return nil
}

func (c *Chain) handleMsg(data []byte) ([]byte, error) {
	if data == nil || len(data) <= 0 {
		log.Debug("Got nil message")
		return nil, errNoMsgContent
	}

	msg := &blockchain.State{}

	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp, err := c.state.diff(msg)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return resp, nil
}

func (c *Chain) handleResponse(resp []byte) {
	if resp == nil || len(resp) == 0 {
		log.Debug("Got nil response")
		return
	}
	msg := &blockchain.StateResponse{}

	err := proto.Unmarshal(resp, msg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//c.log.Debug.Println("Entries in response: ", len(msg.GetMissingEntries()))

	c.state.mergeResponse(msg)

	c.updateState()
}

func (c *Chain) forkResolver() {
	for {
		select {
		case <-c.exitChan:
			return
		case <-time.After(c.forkCheckTimeout):
			if c.state.isForked() {
				log.Info("Sending fork request")
				id := c.state.randomPeerInMajority()

				msg, err := c.fullChain()
				if err != nil {
					log.Error(err.Error())
					continue
				}

				ch, err := c.ifrit.SendToId(id, msg)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				resp := <-ch

				c.mergeChains(resp)
			}
		}
	}
}

func (c *Chain) fullChain() ([]byte, error) {
	c.blockMapMutex.RLock()
	defer c.blockMapMutex.RUnlock()

	state := &blockchain.ChainState{
		Blocks: make(map[uint64]*blockchain.BlockHeader),
	}

	for _, b := range c.blocks {
		header := &blockchain.BlockHeader{
			RootHash: b.getRootHash(),
			PrevHash: b.getPrevHash(),
			BlockNum: b.num,
		}
		state.Blocks[b.num] = header
	}

	return proto.Marshal(state)
}

func (c *Chain) handleForkMsg(msg []byte) ([]byte, error) {
	c.blockMapMutex.RLock()
	defer c.blockMapMutex.RUnlock()

	var err error
	var i uint64
	var response []byte

	log.Info("Got fork merge request")

	content := &blockchain.ChainState{}

	err = proto.Unmarshal(msg, content)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	remoteBlocks := content.GetBlocks()
	if remoteBlocks == nil {
		return nil, errNoChainState
	}

	for i = 1; i <= c.currBlock; i++ {
		localBlock := c.blocks[i]

		if b, ok := remoteBlocks[i]; !ok || !localBlock.cmpRootHash(b.GetRootHash()) {
			log.Info("Found fork point", "blocknum", i, "currblock", c.currBlock)
			response, err = c.partialChain(i)
			if err != nil {
				log.Error(err.Error())
				return nil, err
			}
			break
		}
	}

	return response, nil
}

func (c *Chain) partialChain(start uint64) ([]byte, error) {
	ret := &blockchain.ChainContent{}

	for i := start; i <= c.currBlock; i++ {
		if b, ok := c.blocks[i]; !ok {
			log.Error("Block not present in chain")
			continue
		} else {
			content := &blockchain.BlockContent{
				RootHash: b.getRootHash(),
				PrevHash: b.getPrevHash(),
				Content:  b.content(),
				BlockNum: b.num,
			}
			ret.Blocks = append(ret.Blocks, content)
		}
	}

	return proto.Marshal(ret)
}

func (c *Chain) mergeChains(msg []byte) {
	c.blockMapMutex.Lock()
	defer c.blockMapMutex.Unlock()

	var blocks []*block
	content := &blockchain.ChainContent{}

	err := proto.Unmarshal(msg, content)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, b := range content.GetBlocks() {
		newBlock, err := createBlockWithContent(b.GetPrevHash(), b.GetContent())
		if err != nil {
			log.Error(err.Error())
			return
		}

		newBlock.num = b.GetBlockNum()

		blocks = append(blocks, newBlock)
	}

	prev := c.replaceChain(blocks)

	c.state.changePrev(prev)
}

func (c *Chain) blockLoop() {
	for {
		select {
		case <-c.exitChan:
			return

		case <-time.After(c.blockPeriod):
			c.pickFavouriteBlock()
			c.updateState()
		}
	}
}

func (c *Chain) pickFavouriteBlock() {
	newBlock, expData := c.state.newRound()

	if doExp := c.experiment; doExp {
		c.sendResults(expData)
	}

	// TODO if we get error here, we need a fallback strat
	err := c.addBlock(newBlock)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("New block", "rootHash", string(newBlock.getRootHash()), "prevhash", string(newBlock.getPrevHash()), "number", c.currBlock)
}

func (c *Chain) updateState() {
	bytes, err := c.state.bytes()
	if err != nil {
		log.Error(err.Error())
		return
	}
	c.ifrit.SetGossipContent(bytes)
}

func (c *Chain) replaceChain(blocks []*block) []byte {
	for _, b := range blocks {
		if currBlock, ok := c.blocks[b.num]; !ok {
			c.state.removeEntries(currBlock)
		} else {
			log.Error("Did not have block", "number", b.num)
		}

		c.blocks[b.num] = b

		if b.num > c.currBlock {
			c.currBlock = b.num
		}
	}

	return c.blocks[c.currBlock].getRootHash()
}

func (c *Chain) addBlock(b *block) error {
	c.blockMapMutex.Lock()
	defer c.blockMapMutex.Unlock()

	if b == nil {
		return errNilBlock
	}

	c.currBlock++
	b.num = c.currBlock
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

func (c *Chain) isExp() bool {
	c.expMutex.Lock()
	defer c.expMutex.Unlock()

	if c.expStarted {
		return true
	} else {
		c.expStarted = true
		return false
	}
}

func hashBytes(data []byte) []byte {
	return sha256.New().Sum(data)
}
