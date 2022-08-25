package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	network, err := New()
	fmt.Println(err)
	go handlesig(network)
	network.Start()
}

func handlesig(n *Econetwork) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	for range c {
		fmt.Println("hey")
		n.Stop()
		os.Exit(0)
	}
}
