package main

import (
	"emailchain/network"
	"flag"
)

func main() {

	serverPort := flag.String("port", "9090", "http port number where server will run")
	flag.Parse()
	node := network.NewNode(*serverPort)
	defer node.Wait()
}
