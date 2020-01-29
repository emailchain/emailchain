package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func discoverPeers(n *Network, p Peer){
	log.Printf("Discovering Peers")
	var (
		boot = fmt.Sprintf("http://%s:%s%s", p.Address, p.Port, NODE_PEERS)
	)
	req, err := http.NewRequest(http.MethodGet, boot, nil)
	if err != nil {
		log.Printf("Error sending connect 1: %s", err)
		return
	}
	res, err := n.client.Do(req)
	if err != nil {
		log.Printf("Error sending connect 2: %s", err)
		return
	}
	defer res.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error sending connect 3: %s", err)
		return
	}
	var result map[string][]Peer
	json.Unmarshal(body,&result)

	var newList []Peer
	for _,peer := range result["nodes"]{
		newList = append(newList,peer)
	}
	newList = append(newList, p)
	n.UpdatePeerList(newList)
	for _, host := range result["nodes"] {
		n.node.join <- Peer{host.Port, host.Address}


		var (
			req *http.Request
			res *http.Response
			uri = fmt.Sprintf("http://%s:%s%s", host.Address, host.Port, NODE_PEERS)
		)
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			log.Printf("Error sending connect 4: %s", err)
		}
		res, err = n.client.Do(req)
		if err != nil {
			log.Printf("Error sending connect 5: %s", err)
		}
		defer res.Body.Close()

		var body []byte
		if body, err = ioutil.ReadAll(res.Body); err != nil {
			log.Printf("Error sending connect 6: %s", err)
			return
		}
		var result2 map[string][]Peer
		json.Unmarshal(body,&result2)

		for _, h := range result2["nodes"] {
			n.node.join <- Peer{h.Port, h.Address}
		}
	}
	n.node.join <- p
}
