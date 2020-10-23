package blockchain

import (
	"emailchain/presistence"
)

type Iterator interface {
	Next() *Block
	HasNext() bool
}

type BlockchainIterator struct {
	CurrentHash string
	Index       int64
	db          *presistence.Database
}

func NewBlockchainIterator(bc Blockchain) *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: string(bc.Db.Tip()),
		Index:       bc.LastBlock().Height,
		db:          bc.Db,
	}
}

func (bi *BlockchainIterator) Next() *Block {
	block := DeserializeBlock(bi.db.GetBlock(bi.CurrentHash))
	bi.CurrentHash = block.PreviousHash
	bi.Index--
	return block
}
func (bi *BlockchainIterator) HasNext() bool {
	if bi.Index < 0 {
		return false
	}
	return true
}
