package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/TorchedSammy/Econode/econetwork"
)

func main() {
	network, err := econetwork.New()
	fmt.Println(err)
	go handlesig(network)
	network.Start()
}

func handlesig(n *econetwork.Econetwork) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	for range c {
		fmt.Println("hey")
		n.Stop()
		os.Exit(0)
	}
}
