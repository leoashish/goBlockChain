package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 5000, "TCP Port for blockchain Server")
	flag.Parse()

	app := NewBlockChainServer(uint64(*port))
	app.Run()
}
