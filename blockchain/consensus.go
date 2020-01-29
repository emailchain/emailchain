package blockchain

import (
	"crypto/sha256"
	"math"
	"math/big"
)
const TARGET_BITS uint = 10

type PoW struct {
	Block      Block
	TargetBits uint
	Nonce      int64
}

func NewPow(block Block, targetBits uint, nonce int64) *PoW {
	return &PoW{Block: block, TargetBits: targetBits, Nonce: nonce}
}
// newTarget returns a target number with required leading zero.
func newTarget(targetBits uint) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)
	return target
}
// run the proof of work algorithm
func (pow *PoW) ComputeProof(){
	nonce := int64(0)
	target := newTarget(pow.TargetBits)
	maxNonce := int64(math.MaxInt64)
	for nonce <= maxNonce{
				var hashInt big.Int
		hashInt.SetBytes(Hash(pow.Block.Serialize()))
		if hashInt.Cmp(target) == -1 {
			pow.Nonce = nonce
			pow.Block.Nonce = nonce
			break
		}
		nonce++
		pow.Block.Nonce = nonce
	}
}
// hash data into byte array
func Hash(data []byte) []byte {
	hash := sha256.Sum256(append(data))
	return hash[:]
}
// validate that the proof of work is valid
func IsProofValid(block Block, targetBits uint) bool {
	var hashInt big.Int
	hashInt.SetBytes(Hash(block.Serialize()))
	target := newTarget(targetBits)
	return hashInt.Cmp(target) == -1
}