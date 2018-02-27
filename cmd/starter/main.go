package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	mrand "math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/joonnna/blocks"
	"github.com/joonnna/ifrit"
	"github.com/joonnna/ifrit/bootstrap"
	"github.com/joonnna/ifrit/worm"

	log "github.com/inconshreveable/log15"
)

var (
	clients []string
)

func createClients(requestChan chan interface{}, exitChan chan bool, viz string) {
	for {
		select {
		case <-requestChan:
			var addrs []string
			if length := len(clients); length > 0 {
				idx := mrand.Int() % length
				addrs = append(addrs, clients[idx])
			}

			conf := &ifrit.Config{
				Visualizer: true,
				VisAddr:    viz,
				EntryAddrs: addrs,
			}

			c, err := blocks.NewClient(conf)
			if err != nil {
				fmt.Println(err)
				continue
			}

			clients = append(clients, c.Addr())

			go addPeriodically(c)

			requestChan <- c
		case <-exitChan:
			return
		}

	}
}

func addPeriodically(c *blocks.Client) {
	for {
		time.Sleep(time.Second * 10)
		buf := make([]byte, 50)
		_, err := rand.Read(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		c.Add(buf)
	}
}

func main() {
	var numRings uint
	var vizAddr string
	var wormInterval uint

	runtime.GOMAXPROCS(runtime.NumCPU())

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.UintVar(&numRings, "numRings", 3, "Number of gossip rings to be used")
	args.StringVar(&vizAddr, "vizAddr", "127.0.0.1:8095", "Address of the visualizer(ip:port).")
	args.UintVar(&wormInterval, "wormInterval", 90, "Interval to pull states")
	args.Parse(os.Args[1:])

	ch := make(chan interface{})
	exitChan := make(chan bool)

	r := log.Root()

	h := log.CallerFileHandler(log.Must.FileHandler("/var/log/chainClient", log.TerminalFormat()))

	r.SetHandler(h)

	w := worm.NewWorm(blocks.CmpStates, "hosts", "state", vizAddr, wormInterval)

	l, err := bootstrap.NewLauncher(uint32(numRings), ch, w)
	if err != nil {
		panic(err)
	}

	go createClients(ch, exitChan, vizAddr)
	go l.Start()

	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-channel

	l.ShutDown()
	close(exitChan)

	/*
		var c1 []merkletree.Content
		var c2 []merkletree.Content

		elem1 := &entry{data: []byte("satan i gatan")}
		elem2 := &entry{data: []byte("smeg satan")}
		elem3 := &entry{data: []byte("drittspill")}

		c1 = append(c1, elem1)
		c1 = append(c1, elem2)
		c1 = append(c1, elem3)

		c2 = append(c2, elem3)
		c2 = append(c2, elem2)
		c2 = append(c2, elem1)

		t1, _ := merkletree.NewTree(c1)
		t2, _ := merkletree.NewTree(c2)

		fmt.Println(t1.MerkleRoot())
		fmt.Println(t2.MerkleRoot())

		fmt.Println(bytes.Equal(t1.MerkleRoot(), t2.MerkleRoot()))
	*/
}
