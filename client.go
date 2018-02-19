package blocks

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/joonnna/blocks/core"
	"github.com/joonnna/ifrit/logger"
)

type Client struct {
	log *logger.Log
	ch  *core.Chain

	exitChan chan bool
}

var (
	errInvalidData = errors.New("given data is empty or nil")
)

func NewClient(entryAddr string) (*Client, error) {
	//logFile, err := os.OpenFile("/var/tmp/chainClient")
	logFile, err := os.OpenFile("/var/tmp/chainClient", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log := logger.CreateLogger(fmt.Sprintf("%d", rand.Int()), logFile)

	c, err := core.NewChain(entryAddr, log)
	if err != nil {
		return nil, err
	}

	return &Client{
		ch:       c,
		exitChan: make(chan bool),
		log:      log,
	}, nil
}

func (c *Client) Start() {
	c.log.Info.Println("Started chain client!")
	c.ch.Start()
}

func (c *Client) ShutDown() {
	close(c.exitChan)
}

func (c *Client) Add(data []byte) error {
	if data == nil || len(data) == 0 {
		return errInvalidData
	}
	return c.ch.Add(data)
}
