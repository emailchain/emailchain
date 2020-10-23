package network

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Get a list of peers an node is connected to
func (h *handler) GetPeerList(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodGet {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}
	log.Println("Peer list Requested")
	resp := map[string]interface{}{
		"nodes": h.network.node.NodePeers,
	}
	return Response{resp, http.StatusOK, nil}
}

// Get a list mails in the nodes mempool
func (h *handler) GetPendingMails(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodGet {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}
	log.Println("Pending Mails Requested")
	resp := map[string]interface{}{
		"memPool": h.network.blockchain.Mempool,
	}
	return Response{resp, http.StatusOK, nil}
}

func (h *handler) ConnectNode(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}

	log.Println("Adding node to the blockchain")

	var body map[string][]Peer
	err := json.NewDecoder(r.Body).Decode(&body)

	for _, peer := range body["nodes"] {
		h.network.ConnectToNetwork(peer)
	}
	height := h.network.UpdateChain()

	resp := map[string]interface{}{
		"message":      "New nodes have been added",
		"chain_height": height,
	}

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("fail to register nodes")
		log.Printf("there was an error when trying to register a new node %v\n", err)
	}

	return Response{resp, status, err}
}

func (h *handler) UpdatePeerList(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}

	log.Println("Updating Peer List")

	var body map[string][]Peer
	err := json.NewDecoder(r.Body).Decode(&body)
	go h.network.AddToPeerList(body["nodes"])
	resp := map[string]interface{}{
		"message": "New nodes have been added",
		"nodes":   &h.network.node.NodePeers,
	}

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("fail to register nodes")
		log.Printf("there was an error when trying to register a new node %v\n", err)
	}

	return Response{resp, status, err}
}
