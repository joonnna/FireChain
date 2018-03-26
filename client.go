package blocks

import (
	"errors"

	"github.com/joonnna/blocks/core"
	"github.com/joonnna/ifrit"

	log "github.com/inconshreveable/log15"
)

type Client struct {
	ch *core.Chain
}

var (
	errInvalidData = errors.New("given data is empty or nil")
)

func NewClient(conf *ifrit.Config, hosts, period uint32, expAddr string) (*Client, error) {
	c, err := core.NewChain(conf, hosts, period, expAddr)
	if err != nil {
		return nil, err
	}

	return &Client{
		ch: c,
	}, nil
}

func (c *Client) Start() {
	log.Info("Started chain client")
	c.ch.Start()
}

func (c *Client) ShutDown() {
	log.Info("Stopping chain client")
	c.ch.Stop()
}

func (c *Client) Addr() string {
	return c.ch.Addr()
}

func (c *Client) HttpAddr() string {
	return c.ch.HttpAddr()
}

func (c *Client) Add(data []byte) error {
	if data == nil || len(data) == 0 {
		return errInvalidData
	}
	return c.ch.Add(data)
}

func (c *Client) WaitExp() {
	<-c.ch.ExpChan
}

func CmpStates(first, second []byte) (uint64, error) {
	return core.CmpStates(first, second)
}
