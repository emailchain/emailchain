package blockchain

import (
	. "emailchain/presistence"
	"fmt"
	"log"
	"time"
)

type Blockchain struct {
	Mempool map[string]Email
	Db      *Database
}

func NewBlockchain(nodeID string) *Blockchain {
	database := NewDatabase(nodeID)
	newBlockchain := &Blockchain{
		Mempool: make(map[string]Email),
		Db:      database,
	}
	// Initial generation block
	if string(newBlockchain.Db.Tip()) == "" {
		log.Println("Creating genesis block")
		newBlockchain.GenesisBlock()
	}
	return newBlockchain
}

// add block to chain after making sure it is valid
func (bc *Blockchain) AddBlockToDB(newBlock Block) error {
	if !bc.ValidChain(newBlock) {
		return fmt.Errorf("potential fork on the blockchain")
	}
	if bc.ValidPoW(newBlock) {
		log.Println("Adding valid block to the blockchain")
		newBlockHash := ComputeHashForBlock(newBlock)
		// add block to database
		bc.Db.AddBlock(newBlock.Serialize(), newBlockHash)
		// change tip to the new block
		bc.Db.AddTip(newBlockHash)
		bc.removeFromMempool(newBlock)
	}
	return nil
}
func (bc *Blockchain) GenerateBlock() Block {
	prevBlock := bc.LastBlock()
	prevHash := ComputeHashForBlock(prevBlock)
	var emails []Email
	for _, email := range bc.Mempool {
		emails = append(emails, email)
	}
	newBlock := Block{
		Height:       bc.LastBlock().Height + 1,
		Timestamp:    time.Now().UnixNano(),
		Emails:       emails,
		PreviousHash: prevHash,
	}
	// run proof of work on the block
	pow := NewPow(newBlock, TARGET_BITS, 0)
	pow.ComputeProof()
	newBlock.Nonce = pow.Nonce
	return newBlock
}
func (bc *Blockchain) removeFromMempool(block Block) {
	for _, email := range block.Emails {
		delete(bc.Mempool, email.ID)
	}
}

func (bc *Blockchain) GenesisBlock() Block {
	newBlock := Block{
		Height:       0,
		Timestamp:    0,
		Emails:       nil,
		PreviousHash: "0",
	}
	newBlockHash := ComputeHashForBlock(newBlock)
	// add block to database
	bc.Db.AddBlock(newBlock.Serialize(), newBlockHash)
	// change tip to the new block
	bc.Db.AddTip(newBlockHash)
	return newBlock
}
func (bc *Blockchain) NewEmail(email Email) {
	email.GenerateID()
	bc.Mempool[email.ID] = email
}

func (bc *Blockchain) LastBlock() Block {
	return *DeserializeBlock(bc.Db.GetBlock(string(bc.Db.Tip())))
}
func (bc *Blockchain) GetBlock(hash string) Block {
	return *DeserializeBlock(bc.Db.GetBlock(hash))
}
func (bc *Blockchain) AllBlocks() []Block {
	iterator := NewBlockchainIterator(*bc)
	var listOfBlocks []Block
	for iterator.HasNext() {
		listOfBlocks = append(listOfBlocks, *iterator.Next())
	}
	return listOfBlocks
}
func (bc *Blockchain) AllBlocksFrom(height int64) []Block {
	iterator := NewBlockchainIterator(*bc)
	var listOfBlocks []Block
	for iterator.HasNext() {
		nextBlock := iterator.Next()
		if nextBlock.Height > height {
			listOfBlocks = append(listOfBlocks, *nextBlock)
		}
	}
	return listOfBlocks
}

func (bc *Blockchain) MailBox(pubKey string) []Email {
	iterator := NewBlockchainIterator(*bc)
	var listOfEmails []Email
	for iterator.HasNext() {
		emailsInBlock := iterator.Next().Emails
		for _, email := range emailsInBlock {
			if email.Recipient == pubKey {
				listOfEmails = append(listOfEmails, email)
			}
		}
	}
	return listOfEmails
}
func (bc *Blockchain) Sent(pubKey string) []Email {
	iterator := NewBlockchainIterator(*bc)
	var listOfEmails []Email
	for iterator.HasNext() {
		emailsInBlock := iterator.Next().Emails
		for _, email := range emailsInBlock {
			if email.Sender == pubKey {
				listOfEmails = append(listOfEmails, email)
			}
		}
	}
	return listOfEmails
}

// check if the block has a valid proof of work
func (bc *Blockchain) ValidPoW(block Block) bool {
	if !IsProofValid(block, TARGET_BITS) {
		log.Println("Block Validation Error: Proof of work not valid..")
		return false
	}

	return true
}

// check if the block breaks the chain
func (bc *Blockchain) ValidChain(block Block) bool {
	if ComputeHashForBlock(bc.LastBlock()) != block.PreviousHash {
		log.Println("Block Validation Error: Block breaks chain, possible fork issue. The node needs to be synced")
		return false
	}
	return true
}

// append new blocks to the chain
func (bc *Blockchain) UpdateChain(blocks []Block) error {
	var err error = nil
	for i := len(blocks) - 1; i >= 0; i-- {
		err = bc.AddBlockToDB(blocks[i])
		if err != nil {
			break
		}
	}
	return err
}

// resolve fork
func (bc *Blockchain) ResolveFork(blocks []Block) {
	log.Println("Resolving Fork...")
	// TODO make it more efficent
	bc.GenesisBlock()
	bc.UpdateChain(blocks)

}

// calculate the total work on chain
func (bc *Blockchain) TotalWork() int64 {
	totalWork := int64(0)
	iterator := NewBlockchainIterator(*bc)
	for iterator.HasNext() {
		totalWork += iterator.Next().Nonce
	}
	return totalWork
}
