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
	hashInput := append(email.serialize())
	hash := utils.ComputeHashSha256(hashInput)
	email.ID = hash
}
// Serialize email to slices
func (email *Email) serialize() []byte {
	encoded, err := json.Marshal(email)
	if err != nil {
		log.Panic(err)
	}
	return encoded
}