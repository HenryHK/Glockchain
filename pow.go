package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const targerBits int = 24

// ProofOfWork defines the desired proof of work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork is to generate a target for PoW
// our target here is fixed, to generate a hash with 24 leading 0s in bits
// because we use big int as the target, a successful target can be considered as genrating a number smaller than 1<<(256-24)
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// left shift
	target.Lsh(target, uint(256-targerBits))
	// create a ProofOfWork instance conataining target and the original block
	pow := &ProofOfWork{b, target}
	return pow
}

// prepareData simply combine a block with target and nonce
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targerBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run defines the procedure to work out a valid answer
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for {
		// Step1: prepare the data
		data := pow.prepareData(nonce)
		// Step2: generate sha-256 hash of data
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		// Step3: convert the generated hash to a big int
		hashInt.SetBytes(hash[:])
		// Step4: compare generated int with target
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n\n")
	return nonce, hash[:]
}

// Validate validates block's PoW
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
