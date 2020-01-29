package network

import (
	. "emailchain/blockchain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// get mailbox for public key
func (h *handler) GetMailBox(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}

	log.Println("Getting Mailbox for public key")

	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)

	mailBox := h.network.blockchain.MailBox(body["pubkey"])

	resp := map[string]interface{}{
		"mailBox":   mailBox,
	}

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("failed to get mailbox")
		log.Printf("there was an error : %v\n", err)
	}

	return Response{resp, status, err}
}
// get sent mails for public key
func (h *handler) GetSent(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}

	log.Println("Getting Sent mail for public key")

	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)

	mailBox := h.network.blockchain.Sent(body["pubkey"])

	resp := map[string]interface{}{
		"mailBox":   mailBox,
	}

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("failed to get mailbox")
		log.Printf("there was an error : %v\n", err)
	}

	return Response{resp, status, err}
}

func (h *handler) AddEmail(w io.Writer, r *http.Request) Response {
	if r.Method != http.MethodPost {
		return Response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", r.Method),
		}
	}

	log.Println("Adding Email to Memory pool")
	var email Email
	err := json.NewDecoder(r.Body).Decode(&email)
	h.network.blockchain.NewEmail(email)
	if  len(r.Header[HDR_BDCAST]) == 0 || r.Header[HDR_BDCAST][0] != "true"{
		var value []byte
		value, err = json.Marshal(email)
		message := Message{endpoint: EMAIL_NEW,value:value}
		h.network.node.broadcast <- message
	}


	resp := map[string]string{
		"message": fmt.Sprintf("Email will be added to the chain soon"),
	}

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError
		log.Printf("there was an error when trying to add a Email %v\n", err)
		err = fmt.Errorf("failed to add Email to the blockchain")
	}

	return Response{resp, status, err}
}

