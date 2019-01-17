package main

import (
	"github.com/boltdb/bolt"
	"log"
)

// BlockchainIterator stores the information to iterate persistent DB
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns the block it pointed to and moves pointer to the next block
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic("Error going next using iterator:", err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}
