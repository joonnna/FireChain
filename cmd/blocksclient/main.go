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
	var caAddr, entry, vizAddr, expAddr string
	var logging bool
	var hosts, blockPeriod, blockSize, entrySize uint
	var forkLimit float64

	runtime.GOMAXPROCS(runtime.NumCPU())

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.StringVar(&caAddr, "ca", "", "address(ip:port) of certificate authority")
	args.StringVar(&entry, "entry", "", "address(ip:port) of existing clients")
	args.StringVar(&vizAddr, "viz", "", "address(ip:port) of visualizer")
	args.StringVar(&expAddr, "exp", "", "address of where to send experiment results")
	args.BoolVar(&logging, "log", false, "Bool deciding whether to log")
	args.UintVar(&hosts, "hosts", 0, "How many participants in experiment")
	args.UintVar(&blockPeriod, "btime", 10, "Period between block chosing")
	args.Float64Var(&forkLimit, "flimit", 0.60, "Limit for forking")
	args.UintVar(&blockSize, "blocksize", 25, "Block size in KB")
	args.UintVar(&entrySize, "entrysize", 1024, "Entry size in bytes")
	args.Parse(os.Args[1:])

	r := log.Root()

	if logging {
		h := log.CallerFileHandler(log.Must.FileHandler("blockslog", log.TerminalFormat()))
		r.SetHandler(h)
	} else {
		r.SetHandler(log.DiscardHandler())
	}

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

	c, err := blocks.NewClient(conf, forkLimit, uint32(blockPeriod), uint32(blockSize), uint32(entrySize), expAddr)
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
