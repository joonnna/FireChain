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

func addPeriodically(c *blocks.Client) {
	for {
		time.Sleep(time.Second * 100)
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
	var caAddr string
	var entry string
	var vizAddr string
	var logging bool

	runtime.GOMAXPROCS(runtime.NumCPU())

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.StringVar(&caAddr, "ca", "", "address(ip:port) of certificate authority")
	args.StringVar(&entry, "entry", "", "address(ip:port) of existing clients")
	args.StringVar(&vizAddr, "viz", "", "address(ip:port) of visualizer")
	args.BoolVar(&logging, "log", true, "Bool deciding whether to log")
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

	c, err := blocks.NewClient(conf)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	go c.Start()
	go addPeriodically(c)

	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-channel

	c.ShutDown()
}
