package blockchain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/joonnna/ifrit"
	"github.com/joonnna/ifrit/logger"
	"github.com/joonnna/ifrit/protobuf"
)

type Client struct {
	ifrit *ifrit.Client
	log   *logger.Log
	ch    *chain

	exitChan chan bool
}

func NewClient(entryAddr string) (*Client, error) {
	i, err := ifrit.NewClient(entryAddr, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		ifrit:    i,
		ch:       newChain(),
		exitChan: make(chan bool),
		log:      logger.CreateLogger(fmt.Sprintf("%d", rand.Int()), "chainClient"),
	}, nil
}

func (c *Client) Start() {
	c.log.Info.Println("Started chain client!")
	c.ifrit.RegisterMsgHandler(c.handleMsg)
	go c.ifrit.Start()

	for {
		c.log.Debug.Println("Got here")
		select {
		case <-c.exitChan:
			c.log.Debug.Println("Exited")
			return

		case <-time.After(time.Second * 20):
			err := c.ch.Add(rand.Int31())
			if err != nil {
				c.log.Err.Println(err)
				continue
			}

			buf, err := proto.Marshal(&gossip.Test{Nums: c.ch.State()})
			if err != nil {
				c.log.Err.Println(err)
				continue
			}
			c.log.Debug.Println(len(buf))
			err = c.ifrit.SetGossipContent(buf)
			if err != nil {
				c.log.Err.Println(err)
				continue
			}
		}
	}
}

func (c *Client) handleMsg(data []byte) ([]byte, error) {
	msg := &gossip.Test{}

	err := proto.Unmarshal(data, msg)
	if err != nil {
		c.log.Info.Println(err)
		return nil, err
	}
	//c.log.Info.Println("Received: ", msg.Nums)

	for _, v := range msg.Nums {
		c.ch.Add(v)
	}

	//c.log.Info.Println("Now have: ", c.ch.State())

	return nil, nil
}

func (c *Client) ShutDown() {
	c.ifrit.ShutDown()
	close(c.exitChan)
}

/*
func (c *Client) Add(data []byte) error {
	return c.ch.Add(data)
}
*/
