package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "net/http/pprof"

	"github.com/joonnna/blocks/blockchain"
	"github.com/joonnna/ifrit/bootstrap"
)

func createClients(requestChan chan interface{}, exitChan chan bool, arg string) {
	for {
		select {
		case <-requestChan:
			c, err := blockchain.NewClient(arg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			requestChan <- c
		case <-exitChan:
			return
		}

	}
}

func main() {
	var numRings uint
	var cpuprofile, memprofile string

	runtime.GOMAXPROCS(runtime.NumCPU())

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.UintVar(&numRings, "numRings", 3, "Number of gossip rings to be used")

	args.Parse(os.Args[1:])

	ch := make(chan interface{})
	exitChan := make(chan bool)

	l, err := bootstrap.NewLauncher(uint32(numRings), ch)
	if err != nil {
		panic(err)
	}

	go createClients(ch, exitChan, l.EntryAddr)
	go l.Start()

	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-channel

	l.ShutDown()
	close(exitChan)
}
