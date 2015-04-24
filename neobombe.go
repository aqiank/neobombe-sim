package main

import (
	"os"
	"os/signal"
	"sync"

	"github.com/jackyb/neobombe-sim/osc"
	"github.com/jackyb/neobombe-sim/server"
	"github.com/jackyb/neobombe-sim/simulation"
	"github.com/jackyb/neobombe-sim/twitter"
)

var sigChan = make(chan os.Signal, 1)

var channels = map[string]interface{}{
	"signal":    sigChan,
	"message":   make(chan string, 1),
	"encrypted": make(chan string, 1),
	"decrypted": make(chan string, 1),
	"state":     make(chan simulation.Bombe, 1),
	"osc":       make(chan simulation.Bombe, 1),
}

func main() {
	wg := sync.WaitGroup{}

	signal.Notify(sigChan, os.Interrupt, os.Kill)

	waitFunc(&wg, twitter.Run)
	waitFunc(&wg, server.Run)
	waitFunc(&wg, simulation.Run)
	waitFunc(&wg, osc.Run)

	wg.Wait()
}

func waitFunc(wg *sync.WaitGroup, fn func(cs map[string]interface{})) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fn(channels)
	}()
}
