package network

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type handler struct {
	network *Network
}
type Response struct {
	value      interface{}
	statusCode int
	err        error
}

func NewHandler(network *Network) http.Handler {

	h := handler{network: network}

	mux := http.NewServeMux()
	mux.HandleFunc(NODE_CONNECT, BuildHTTPResponse(h.ConnectNode))
	mux.HandleFunc(NODE_PEERS, BuildHTTPResponse(h.GetPeerList))
	mux.HandleFunc(NODE_PEERS_UPDATE, BuildHTTPResponse(h.UpdatePeerList))
	mux.HandleFunc(NODE_PENDINGMAILS, BuildHTTPResponse(h.GetPendingMails))
	mux.HandleFunc(EMAIL_NEW, BuildHTTPResponse(h.AddEmail))
	mux.HandleFunc(EMAIL_MAILBOX, BuildHTTPResponse(h.GetMailBox))
	mux.HandleFunc(EMAIL_SENT, BuildHTTPResponse(h.GetSent))
	mux.HandleFunc(CHAIN_GENERATE, BuildHTTPResponse(h.GenerateBlock))
	mux.HandleFunc(CHAIN_VIEW, BuildHTTPResponse(h.ViewBlockchain))
	mux.HandleFunc(CHAIN_INFO, BuildHTTPResponse(h.GetBlockchainInfo))
	mux.HandleFunc(CHAIN_GETBLOCK, BuildHTTPResponse(h.GetBlockByHash))
	mux.HandleFunc(CHAIN_SYNC, BuildHTTPResponse(h.SyncChain))
	mux.HandleFunc(CHAIN_MISSINGBLOCKS, BuildHTTPResponse(h.GetMissingBlocks))
	mux.HandleFunc(CHAIN_ADD, BuildHTTPResponse(h.AddBlock))
	return mux
}

func BuildHTTPResponse(h func(io.Writer, *http.Request) Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := h(w, r)
		msg := resp.value
		if resp.err != nil {
			msg = resp.err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.statusCode)
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			log.Printf("could not encode Response to output: %v", err)
		}
	}
}
