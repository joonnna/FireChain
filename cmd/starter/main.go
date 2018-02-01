package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	_ "net/http/pprof"

	"github.com/joonnna/blocks/blockchain"
	"github.com/joonnna/ifrit/ifrit"
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
	args.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile `file`")
	args.StringVar(&memprofile, "memprofile", "", "write memory profile to `file`")

	args.Parse(os.Args[1:])

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

	ch := make(chan interface{})
	exitChan := make(chan bool)

	l, err := ifrit.NewLauncher(uint32(numRings), ch)
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
