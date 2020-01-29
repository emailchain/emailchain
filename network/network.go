package network

import (
	"bytes"
	. "emailchain/blockchain"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/rs/cors"
)

type Network struct {
	node       *Node
	blockchain *Blockchain
	client     *http.Client
}

func NewNetwork(n *Node) *Network {
	// TODO change node ID to something better
	nodeID := n.Self.Address + n.Self.Port
	blockchain := NewBlockchain(nodeID)
	return &Network{
		node:       n,
		blockchain: blockchain,
		client:     &http.Client{},
	}
}

// Starts an HTTP server
func (n *Network) Start() {
	var s = NewHandler(n)
	var host = fmt.Sprintf("%s:%s", n.node.Self.Address, n.node.Self.Port)
	http.Handle("/", s)
	handler := cors.Default().Handler(s)
	http.ListenAndServe(fmt.Sprintf(host), handler)
}

// Discover and connect to peers from a single peer
func (n *Network) ConnectToNetwork(p Peer) {
	discoverPeers(n, p)
}
func (n *Network) Broadcast(message Message) {
	for _, peer := range n.node.NodePeers {
		uri := fmt.Sprintf(HTTP_URL, peer.Address, peer.Port, message.endpoint)
		body := bytes.NewBuffer(message.value)
		req, err := http.NewRequest(http.MethodPost, uri, body)
		if err != nil {
			log.Panic(err)
		}
		req.Header.Set(HDR_BDCAST, "true")
		req.Header.Set(HDR_ADDR, n.node.Self.Address)
		req.Header.Set(HDR_PORT, n.node.Self.Port)
		_, err = n.client.Do(req)
		if err != nil {
			log.Println("Failed to Broadcast Message:"+err.Error())
			// TODO remove peer form peer list
		}
	}
}
func (n *Network) UpdatePeerList(peers []Peer) {
	for _, peer := range peers {
		uri := fmt.Sprintf(HTTP_URL, peer.Address, peer.Port, NODE_PEERS_UPDATE)
		var body = make(map[string][]Peer)
		peers = append(peers, n.node.Self)
		body["nodes"] = peers
		result, err := json.Marshal(body)
		var bodySend = bytes.NewBuffer(result)
		req, err := http.NewRequest(http.MethodPost, uri, bodySend)
		if err != nil {
			log.Panic(err)
		}
		_, err = n.client.Do(req)
		if err != nil {
			log.Panic(err)
		}
	}
}
func (n *Network) AddToPeerList(peers []Peer) {
	for _, peer := range peers {
		n.node.join <- peer
	}
}

// update the blockchain by requesting peers to submit missing blocks and choose the longest one to include
func (n *Network) UpdateChain() int64 {
	log.Println("Updating chain")
	peers := n.node.NodePeers
	for _, peer := range peers {
		height := n.blockchain.LastBlock().Height
		uri := fmt.Sprintf(HTTP_URL, peer.Address, peer.Port, CHAIN_MISSINGBLOCKS)
		var body = make(map[string]int64)
		body["height"] = height
		result, err := json.Marshal(body)
		var bodySend = bytes.NewBuffer(result)
		req, err := http.NewRequest(http.MethodPost, uri, bodySend)
		if err != nil {
			log.Panic(err)
		}
		resp, err := n.client.Do(req)
		if err != nil {
			log.Panic(err)
		}
		var b map[string][]Block
		err = json.NewDecoder(resp.Body).Decode(&b)
		if err != nil {
			log.Panic(err)
		}
		blocks := b["blocks"]
		if blocks != nil {
			err:= n.blockchain.UpdateChain(blocks)
			if err != nil{
				n.ResolveFork()
			}
		}
	}
	return n.blockchain.LastBlock().Height
}

// get missing blocks
func (n *Network) GetMissingBlocks(height int64, peer Peer) []Block{
	uri := fmt.Sprintf(HTTP_URL, peer.Address, peer.Port, CHAIN_MISSINGBLOCKS)
	var body = make(map[string]int64)
	body["height"] = height
	result, err := json.Marshal(body)
	var bodySend = bytes.NewBuffer(result)
	req, err := http.NewRequest(http.MethodPost, uri, bodySend)
	if err != nil {
		log.Panic(err)
	}
	resp, err := n.client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	var b map[string][]Block
	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	blocks := b["blocks"]
	return blocks
}
// find the block with the highest chain
func (n *Network) HighestChain() Peer {
	height := int64(0)
	highestPow := int64(0)
	peers := n.node.NodePeers
	maxHeightPeer := peers[0]
	for _, peer := range peers {
		var (
			req *http.Request
			res *http.Response
			uri = fmt.Sprintf(HTTP_URL, peer.Address, peer.Port, CHAIN_INFO)
		)
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			log.Panic(err)
		}
		res, err = n.client.Do(req)
		if err != nil {
			log.Panic(err)
		}
		defer res.Body.Close()

		var body []byte
		if body, err = ioutil.ReadAll(res.Body); err != nil {
			log.Panic(err)
		}
		var result map[string]int64
		json.Unmarshal(body, &result)
		//log.Println("---start---")
		//log.Println("port: "+ peer.Port)
		//log.Println(result["tip_height"])
		//log.Println(result["total_pow"])
		//log.Println("---end---")
		if height < result["tip_height"] && highestPow < result["total_pow"]{

			height = result["tip_height"]
			highestPow = result["total_pow"]
			maxHeightPeer = peer
		}
	}
	return maxHeightPeer
}
// resolve fork issue
func (n *Network) ResolveFork(){
	// TODO make it more efficient
	peer := n.HighestChain()
	log.Println("highest port: "+peer.Port)
	height := int64(0)
	blocks := n.GetMissingBlocks(height,peer)
	n.blockchain.ResolveFork(blocks)
}
