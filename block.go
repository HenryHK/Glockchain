package main

import (
	"time"
)

// Block defines the structure of blocks in the blockchain
// Timestamp marks the time, Data carries data, PrevBlockHash stores the hash of previsous block, Hash is the hash of the block itself
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock is used to create new block in the block chain
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock create the genesis block(the first block) of the blockchain
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
