package main

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

var maxNonce = math.MaxInt64

// TARGETBITS define the mining difficulty
const TARGETBITS = 4

// ProofOfWork represents a block mined with a target difficulty
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds a ProofOfWork
func NewProofOfWork(block *Block) *ProofOfWork {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	var i, e = big.NewInt(2), big.NewInt(256 - TARGETBITS)
	// does the exponent calculation from the readme
	i.Exp(i, e, nil)
	return &ProofOfWork{block: block, target: i}
}

// setupHeader prepare the header of the block
func (pow *ProofOfWork) setupHeader() []byte {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	var headers []byte
	headers = append(headers, pow.block.PrevBlockHash...)
	headers = append(headers, pow.block.HashTransactions()...)
	convTime := IntToHex(pow.block.Timestamp)
	headers = append(headers, convTime...)
	convTarget := IntToHex(TARGETBITS)
	headers = append(headers, convTarget...)
	return headers
}

// addNonce adds a nonce to the header
func addNonce(nonce int, header []byte) []byte {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	convNonce := IntToHex(int64(nonce))
	header = append(header, convNonce...)
	return header
}

// Run performs the proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	header := pow.setupHeader()
	// iterates nonce until maxNonce
	for nonce := 0; nonce < maxNonce; nonce += 1 {
		// adds nonce to the header
		nonced := addNonce(nonce, header)
		// hashed the header
		hashed := sha256.Sum256(nonced)
		// converts the hashed header to a bigint
		hashInt := new(big.Int).SetBytes(hashed[:])
		// compares the bigint to the target
		if hashInt.CmpAbs(pow.target) == -1 {
			return nonce, hashed[:]
		}
	}
	return 0, nil
}

func (pow *ProofOfWork) RunCompete(nonce int) (int, []byte) {
	// the header is going to be remade a bunch which sucks but it is what it is
	header := pow.setupHeader()
	// pretty much the same as normal run except if just returns if the nonce is wrong, also the nonce comes externally
	nonced := addNonce(nonce, header)
	hashed := sha256.Sum256(nonced)
	hashInt := new(big.Int).SetBytes(hashed[:])
	if hashInt.CmpAbs(pow.target) == -1 {
		return nonce, hashed[:]
	}
	return 0, nil
}

// Validate validates block's Proof-Of-Work
// This function just validates if the block header hash
// is less than the target AND equals to the mined block hash.
func (pow *ProofOfWork) Validate() bool {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	// block header
	header := pow.block.Hash
	// converts the block header to a big int
	headInt := new(big.Int).SetBytes(header)
	// grabs the nonce to reconstruct the header
	nonce := pow.block.Nonce
	// adds the nonce to the pow header
	nonced := addNonce(nonce, pow.setupHeader())
	// hashes the nonce
	hashNonce := sha256.Sum256(nonced)
	// checks if the target is met and the hash is equal
	if headInt.CmpAbs(pow.target) == -1 && bytes.Equal(pow.block.Hash, hashNonce[:]) {
		return true
	}
	return false
}
