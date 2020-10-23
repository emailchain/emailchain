package blockchain

import (
	"emailchain/utils"
	"encoding/json"
	"log"
)

type Email struct {
	ID        string `json:id`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

// Generate a unique ID for emails
func (email *Email) GenerateID() {
	//timeByte := make([]byte,8)
	//binary.LittleEndian.PutUint64(timeByte, uint64(time.Now().UnixNano()))
	hashInput := append(email.Serialize())
	hash := utils.ComputeHashSha256(hashInput)
	email.ID = hash
}
func (email *Email) Serialize() []byte {
	encoded, err := json.Marshal(email)
	if err != nil {
		log.Panic(err)
	}
	return encoded
}
func (email *Email) Deserialize(d []byte) *Email {
	var e Email
	err := json.Unmarshal(d, &e)
	if err != nil {
		log.Panic(err.Error())
	}
	return &e
}
