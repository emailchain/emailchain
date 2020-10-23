package blockchain

import (
	"emailchain/utils"
	"encoding/json"
	"log"
)
// Email structure
type Email struct {
	ID        string `json:id`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

// Generate a unique ID for emails
func (email *Email) GenerateID() {
	hashInput := append(email.Serialize())
	hash := utils.ComputeHashSha256(hashInput)
	email.ID = hash
}
