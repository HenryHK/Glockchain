package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block defines the structure of blocks in the blockchain
// Timestamp marks the time, Data carries data, PrevBlockHash stores the hash of previsous block, Hash is the hash of the block itself
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// SetHash calculates the hash and set it
func (b *Block) SetHash() {
	// convert timestamp in int64 to a string and then convert the result to bytes
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// headers is the combination of all info conatined in the block
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	// calculate the hash using the headers
	hash := sha256.Sum256(headers)
	//set the hash
	b.Hash = hash[:]
}

// NewBlock is used to create new block in the block chain
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}
