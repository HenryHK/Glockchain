package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const blocksBucket = "blocks"
const dbFile = "blockchain.db"

// Blockchain is the chain holding blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock is to add a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	newBlock := NewBlock(data, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic("Error adding block into db:", err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic("Error adding block into db:", err)
		}
		bc.tip = newBlock.Hash
		return nil
	})
}

// NewBlockchain create a new blockchain with the first block is genesis block. if there exists blockchain already,do nothing.
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic("Error opening db file:", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic("Error creating bucket:", err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic("Error put genesis into db:", err)
			}
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic("Error put genesis into db:", err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc
}
