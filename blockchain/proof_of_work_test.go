package blockchain

import (
	"testing"
)

func TestPoW_ComputeProof(t *testing.T) {
	block := Block{
		Height:       0,
		Timestamp:    0,
		Emails:       nil,
		PreviousHash: "0",
	}
	pow := NewPow(block,TARGET_BITS,0)
	pow.ComputeProof()
	if !IsProofValid(pow.Block,TARGET_BITS){
		t.Fail()
	}
}