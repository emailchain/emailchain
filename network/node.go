package network

import (
	"log"
	"strconv"
	"sync"
)
type Node struct {
	Self    Peer
	NodePeers Peers
	network *Network
	broadcast  chan Message
	join       chan Peer
	leave      chan Peer

	waiter *sync.WaitGroup
	connected bool
}

func NewNode(p string) (n *Node) {
	port,err := strconv.Atoi(p)
	if err != nil{
		log.Panic(err)
	}
	n = &Node{
		Self:       Me(port),
		NodePeers:    Peers{},
		broadcast:     make(chan Message),
		join:       make(chan Peer),
		waiter:     &sync.WaitGroup{},
		connected:  true,
	}
	n.startServices()
	return
}

// Keep node alive
func (n *Node) Wait() {
	n.waiter.Wait()
}
// Broadcast message to all peers
func (n *Node) Broadcast(m Message) {
	n.broadcast <- m
}

// start all goroutines
func (n *Node) startServices() {
	n.waiter.Add(1)
	go n.nodeLoop()
	go n.httpServer()
}

// main loop to handle network events
func (n *Node) nodeLoop() {
	log.Printf("Started Main Node Loop ")
	for {
		select {
		case m := <-n.broadcast:
			if len(n.NodePeers) > 0 {
				go n.network.Broadcast(m)
				log.Printf("Broadcasting message")
			} else {
				log.Printf("Broadcasting aborted. Empty network!")
			}
		case p := <-n.join:
			if !n.NodePeers.contains(p) && !n.Self.isMe(p) {
				n.NodePeers = append(n.NodePeers, p)
				log.Printf("Connected to [%s:%s]", p.Address, p.Port)

			}

		case p := <-n.leave:
			if n.NodePeers.contains(p) && !n.Self.isMe(p) {
				n.NodePeers = n.NodePeers.delete(p)
				log.Printf("Disconnected From [%s:%s]", p.Address, p.Port)
			}
		}
	}
}

func (n *Node) httpServer() {
	n.network = NewNetwork(n)
	log.Printf("Listen at %s:%s", n.Self.Address, n.Self.Port)
	n.network.Start()
}

