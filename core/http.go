package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	log "github.com/inconshreveable/log15"
	"github.com/joonnna/blocks/protobuf"
)

func initHttp() (net.Listener, error) {
	l, err := net.Listen("tcp", "localhost:")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return l, nil
}

func (c *Chain) startHttp() {
	r := mux.NewRouter()

	r.HandleFunc("/state", c.stateHandler)
	r.HandleFunc("/hosts", c.hostsHandler)

	err := http.Serve(c.httpListener, r)
	if err != nil {
		log.Error(err.Error())
	}
}

func (c *Chain) stopHttp() {
	c.httpListener.Close()
}

func (c *Chain) stateHandler(w http.ResponseWriter, r *http.Request) {
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	bytes, err := c.getChain()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Error(err.Error())
	}
}

func (c *Chain) hostsHandler(w http.ResponseWriter, r *http.Request) {
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	addrs := c.state.getHttpAddrs()

	for _, a := range addrs {
		_, err := w.Write([]byte(fmt.Sprintf("%s\n", a)))
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (c *Chain) getChain() ([]byte, error) {
	c.blockMapMutex.RLock()
	defer c.blockMapMutex.RUnlock()

	payload := &blockchain.WormPayload{
		Blocks: make(map[uint64]*blockchain.BlockHeader),
	}

	for blockNum, b := range c.blocks {
		payload.Blocks[blockNum] = &blockchain.BlockHeader{
			RootHash: b.getRootHash(),
			PrevHash: b.getPrevHash(),
		}
	}

	return proto.Marshal(payload)
}

func CmpStates(first, second []byte) (uint64, error) {
	var maxBlock uint64
	s1 := &blockchain.WormPayload{}
	s2 := &blockchain.WormPayload{}

	err := proto.Unmarshal(first, s1)
	if err != nil {
		return 0, err
	}

	err = proto.Unmarshal(second, s2)
	if err != nil {
		return 0, err
	}

	m1 := s1.GetBlocks()
	m2 := s2.GetBlocks()

	if m1 == nil && m2 == nil {
		return 0, nil
	} else if m1 == nil || m2 == nil {
		return 0, errors.New("Nil maps")
	}

	for blockNum, b1 := range m1 {
		if blockNum > maxBlock {
			maxBlock = blockNum
		}
		if b2, exists := m2[blockNum]; exists {
			if !bytes.Equal(b1.GetRootHash(), b2.GetRootHash()) {
				return maxBlock, errors.New(fmt.Sprintf("Different root hashes for block %d", blockNum))
			}

			if !bytes.Equal(b1.GetPrevHash(), b2.GetPrevHash()) {
				return maxBlock, errors.New(fmt.Sprintf("Different prev hashes for block %d", blockNum))
			}

		} else {
			return maxBlock, errors.New(fmt.Sprintf("Not equal, both do not have block %d", blockNum))
		}
	}

	return maxBlock, nil
}
