package core

import (
	"crypto/sha256"
	"errors"
	"fmt"
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

func NewChain(conf *ifrit.Config, hosts, blockPeriod uint32, expAddr string) (*Chain, error) {
	i, err := ifrit.NewClient(conf)
	if err != nil {
		return nil, err
	}

	l, err := initHttp()
	if err != nil {
		return nil, err
	}

	ip := netutil.GetLocalIP()
	port := strings.Split(l.Addr().String(), ":")[1]
	httpAddr := fmt.Sprintf("%s:%s", ip, port)

	exp := false
	if expAddr != "" {
		exp = true
	}

	return &Chain{
		blocks:       make(map[uint64]*block),
		state:        newState(i, httpAddr, hosts),
		ifrit:        i,
		httpListener: l,
		exitChan:     make(chan bool),
		blockPeriod:  time.Minute * time.Duration(blockPeriod),
		ExpChan:      make(chan bool),
		experiment:   exp,
		expAddr:      expAddr,
		httpAddr:     httpAddr,
	}, nil
}

func (c *Chain) Start() {
	go c.ifrit.Start()
	go c.startHttp()

	<-c.ExpChan
	c.ifrit.RegisterGossipHandler(c.handleMsg)
	c.ifrit.RegisterResponseHandler(c.handleResponse)
	c.ifrit.RecordGossipRounds()
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
	newBlock := c.state.newRound()

	if doExp := c.experiment; doExp {
		c.sendResults(newBlock.getRootHash())
	}

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
