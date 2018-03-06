package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	_ "net/http/pprof"

	log "github.com/inconshreveable/log15"
	"github.com/joonnna/blocks"
	"github.com/joonnna/ifrit"
)

func main() {
	var caAddr string
	var entry string
	var vizAddr string

	runtime.GOMAXPROCS(runtime.NumCPU())

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.StringVar(&caAddr, "ca", "", "address(ip:port) of certificate authority")
	args.StringVar(&entry, "entry", "", "address(ip:port) of existing clients")
	args.StringVar(&vizAddr, "viz", "", "address(ip:port) of visualizer")
	args.Parse(os.Args[1:])

	r := log.Root()

	h := log.CallerFileHandler(log.Must.FileHandler("blockslog", log.TerminalFormat()))

	r.SetHandler(h)

	log.Info("Starting block client")

	entryAddrs := strings.Split(entry, ",")

	ca := false
	if caAddr != "" {
		ca = true
	}

	viz := false
	if vizAddr != "" {
		viz = true
	}

	conf := &ifrit.Config{
		Ca:         ca,
		CaAddr:     caAddr,
		EntryAddrs: entryAddrs,
		Visualizer: viz,
		VisAddr:    vizAddr,
	}

	c, err := blocks.NewClient(conf)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	go c.Start()

	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-channel

	c.ShutDown()
}
