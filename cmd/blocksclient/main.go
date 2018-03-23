package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "net/http/pprof"

	log "github.com/inconshreveable/log15"
	"github.com/joonnna/blocks"
	"github.com/joonnna/ifrit"
)

func fillFirstBlock(c *blocks.Client) {
	for {
		buf := make([]byte, 50)
		_, err := rand.Read(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = c.Add(buf)
		if err != nil {
			break
		}
	}
}

func addPeriodically(c *blocks.Client) {
	for {
		time.Sleep(time.Second * 50)
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
	var caAddr, entry, vizAddr, expAddr string
	var logging bool
	var saturation, hosts, blockPeriod uint

	runtime.GOMAXPROCS(runtime.NumCPU())

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.StringVar(&caAddr, "ca", "", "address(ip:port) of certificate authority")
	args.StringVar(&entry, "entry", "", "address(ip:port) of existing clients")
	args.StringVar(&vizAddr, "viz", "", "address(ip:port) of visualizer")
	args.StringVar(&expAddr, "exp", "", "address of where to send experiment results")
	args.BoolVar(&logging, "log", false, "Bool deciding whether to log")
	args.UintVar(&saturation, "wait", 0, "Timeout to start creating block content")
	args.UintVar(&hosts, "hosts", 0, "How many participants in experiment")
	args.UintVar(&blockPeriod, "btime", 10, "Period between block chosing")
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

	c, err := blocks.NewClient(conf, uint32(saturation), uint32(hosts), uint32(blockPeriod), expAddr)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	go c.Start()

	time.Sleep(time.Minute * time.Duration(saturation))
	fillFirstBlock(c)
	go addPeriodically(c)

	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-channel

	c.ShutDown()
}
