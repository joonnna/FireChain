package client

import "github.com/joonnna/ifrit/ifrit"

type Client struct {
	ifrit *ifrit.Client

	ch *blockchain.Chain
}

func NewClient(entryAddr string) {
	return &Client{
		ifrit: ifrit.NewClient(entryAddr),
		c:     blockchain.NewChain(),
	}
}

func (c *Client) Add(data []byte) error {
	return c.ch.Add(data)
}
