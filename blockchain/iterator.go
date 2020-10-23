package blockchain

import (
	"emailchain/database"
)
// Iterator interface
type Iterator interface {
	Next() *Block
	HasNext() bool
}
// Blockchain Iterator structure
type BlockchainIterator struct {
	CurrentHash string
	Index       int64
	db          *database.DB
}
// Create new iterator
func NewBlockchainIterator(bc Blockchain) *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: string(bc.Db.Tip()),
		Index:       bc.LastBlock().Height,
		db:          bc.Db,
	}
}
// Go to next block in the blockchain database
func (bi *BlockchainIterator) Next() *Block {
	block := DeserializeBlock(bi.db.GetBlock(bi.CurrentHash))
	bi.CurrentHash = block.PreviousHash
	bi.Index--
	return block
}
// Check if there is a next block in the blockchain database
func (bi *BlockchainIterator) HasNext() bool {
	if bi.Index < 0 {
		return false
	}
	return true
}
