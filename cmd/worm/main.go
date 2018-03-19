package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/joonnna/blocks"
	"github.com/joonnna/ifrit/worm"
)

func main() {
	var vizAddr, entryAddr string
	var wormInterval uint

	runtime.GOMAXPROCS(runtime.NumCPU())

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.StringVar(&vizAddr, "vizAddr", "129.242.19.135:8095", "Address of the visualizer (ip:port)")
	args.StringVar(&entryAddr, "entry", "", "Address of an existing participant")
	args.UintVar(&wormInterval, "wormInterval", 90, "Interval to pull states")

	args.Parse(os.Args[1:])

	exitChan := make(chan bool)

	r := log.Root()

	h := log.CallerFileHandler(log.Must.FileHandler("/var/log/wormlog", log.TerminalFormat()))

	r.SetHandler(h)

	w := worm.NewWorm(blocks.CmpStates, "hosts", "state", vizAddr, wormInterval)
	w.AddHost(entryAddr)
	w.Start()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-channel

	w.Stop()
	close(exitChan)
}
