package blockchain

import (
	"bytes"
	"emailchain/utils"
	"encoding/binary"
	"encoding/json"
	"log"
)

type Block struct {
	Height       int64   `json:"height"`
	Timestamp    int64   `json:"timestamp"`
	Emails       []Email `json:"emails"`
	PreviousHash string  `json:"previous_hash"`
	Nonce        int64 `json:"nonce"`
}

func ComputeHashForBlock(block Block) string {
	var buf bytes.Buffer
	// Data for binary.Write must be a fixed-size value or a slice of fixed-size values,
	// or a pointer to such data.
	jsonblock, marshalErr := json.Marshal(block)
	if marshalErr != nil {
		log.Fatalf("Could not marshal block: %s", marshalErr.Error())
	}
	hashingErr := binary.Write(&buf, binary.BigEndian, jsonblock)
	if hashingErr != nil {
		log.Fatalf("Could not hash block: %s", hashingErr.Error())
	}
	return utils.ComputeHashSha256(buf.Bytes())
}
func (b *Block) Serialize() []byte {
	encoded, err := json.Marshal(b)
	if err != nil {
		log.Panic(err)
	}
	return encoded
}
func DeserializeBlock(d [] byte) *Block {
	var block Block
	err := json.Unmarshal(d, &block)
	if err != nil {
		log.Println(err)
	}
	return &block
}

