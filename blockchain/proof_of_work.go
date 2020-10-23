package blockchain

import (
	"crypto/sha256"
	"math"
	"math/big"
)

const TARGET_BITS uint = 10
// Proof of Work struct
type PoW struct {
	Block      Block
	TargetBits uint
	Nonce      int64
}
// Create a new Proof of Work
func NewPow(block Block, targetBits uint, nonce int64) *PoW {
	return &PoW{Block: block, TargetBits: targetBits, Nonce: nonce}
}

// Returns a target number with required leading zero.
func newTarget(targetBits uint) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)
	return target
}

// Run the proof of work algorithm
func (pow *PoW) ComputeProof() {
	nonce := int64(0)
	target := newTarget(pow.TargetBits)
	maxNonce := int64(math.MaxInt64)
	for nonce <= maxNonce {
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

// Hash data into byte array
func Hash(data []byte) []byte {
	hash := sha256.Sum256(append(data))
	return hash[:]
}

// Validate the Proof of Work
func IsProofValid(block Block, targetBits uint) bool {
	var hashInt big.Int
	hashInt.SetBytes(Hash(block.Serialize()))
	target := newTarget(targetBits)
	return hashInt.Cmp(target) == -1
}
