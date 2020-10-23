package network

import (
	. "emailchain/blockchain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// get block by hash
func (h *handler) GetBlockByHash(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}

	log.Println("Getting Block by hash")

	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)

	block := h.network.blockchain.GetBlock(body["hash"])
	hash := ComputeHashForBlock(block)
	resp := map[string]interface{}{
		"block": block,
		"hash":  hash,
		"valid": IsProofValid(block, TARGET_BITS),
	}

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("fail to get block")
		log.Printf("there was an error:  %v\n", err)
	}

	return Response{resp, status, err}
}

// get missing blocks form current height
func (h *handler) GetMissingBlocks(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}
	log.Println("Getting missing blocks")
	var body map[string]int64
	err := json.NewDecoder(r.Body).Decode(&body)
	blocks := h.network.blockchain.AllBlocksFrom(body["height"])

	resp := map[string]interface{}{
		"blocks": blocks,
	}

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("fail to get block")
		log.Printf("there was an error:  %v\n", err)
	}

	return Response{resp, status, err}
}

// synchronize chain with other peers
func (h *handler) SyncChain(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodGet {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}
	log.Println("Sync Chain requested")
	h.network.UpdateChain()
	resp := map[string]interface{}{
		"tip_height": DeserializeBlock(h.network.blockchain.Db.GetBlock(string(h.network.blockchain.Db.Tip()))).Height,
	}
	return Response{resp, http.StatusOK, nil}
}

// view the whole blockchain
func (h *handler) ViewBlockchain(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodGet {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}
	log.Println("View Blockchain requested")

	resp := map[string]interface{}{
		"chain": h.network.blockchain.AllBlocks(),
	}
	return Response{resp, http.StatusOK, nil}
}

// get blockchain info
func (h *handler) GetBlockchainInfo(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodGet {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}
	log.Println("Blockchain Info requested")

	resp := map[string]interface{}{
		"tip_height": h.network.blockchain.LastBlock().Height,
		"total_pow":  h.network.blockchain.TotalWork(),
	}
	return Response{resp, http.StatusOK, nil}
}

func (h *handler) GenerateBlock(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodGet {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}

	log.Println("Generating Block")

	// Forge the new Block by adding it to the chain
	var err error
	block := h.network.blockchain.GenerateBlock()
	if len(r.Header[HDR_BDCAST]) == 0 || r.Header[HDR_BDCAST][0] != "true" {
		var value []byte
		value, err = json.Marshal(block)
		message := Message{endpoint: CHAIN_ADD, value: value}
		h.network.node.broadcast <- message
		syncError := h.network.blockchain.AddBlockToDB(block)
		if syncError != nil {
			h.network.ResolveFork()
		}
	}
	status := http.StatusOK
	message := "New Block Generated"
	if err != nil {
		status = http.StatusInternalServerError
		log.Printf("there was an error when trying to Generate block %v\n", err)
		err = fmt.Errorf("failed to generate block")
		message = "failed to generate block"
	}
	resp := map[string]interface{}{"message": message, "block": block}
	return Response{resp, status, nil}
}

func (h *handler) AddBlock(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}

	log.Println("Adding Block to blockchain")
	var block Block
	err := json.NewDecoder(r.Body).Decode(&block)
	syncError := h.network.blockchain.AddBlockToDB(block)
	if syncError != nil {
		h.network.ResolveFork()
	}
	resp := map[string]string{
		"message": fmt.Sprintf("New Generated Block added to the blockchain"),
	}

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		log.Printf("there was an error when trying to add a Block %v\n", err)
		err = fmt.Errorf("failed to add Block to the blockchain")
	}

	return Response{resp, status, err}
}
