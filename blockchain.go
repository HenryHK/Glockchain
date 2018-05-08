package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// Blockchain is the chain holding blocks
type Blockchain struct {
	tip []byte   // only the tip is stored which is the hash of the last block in the chain
	db  *bolt.DB // store a DB connection, this should be open and kept while the program is running
}

// BlockchainIterator is the iterator for inspecting the blockchain.
//	This will return the next block from a blockchain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns the next block from a blockchain in trace back order
func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		fmt.Println("Get Block Error:", err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}

// Iterator implements interator interface
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

// AddBlock is to add a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// get the hash of the last block in the chain
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	})

	if err != nil {
		fmt.Println("Get Last Hash ErrorL", err)
	}

	// create new block
	newBlock := NewBlock(data, lastHash)
	// add new block to chain
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		// put new block into db
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			fmt.Println("Add New Block Error:", err)
		}
		// update tip in db
		err = b.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			fmt.Println("Error Update Tip:", err)
		}
		// update the blockchain
		bc.tip = newBlock.Hash
		return nil
	})
}

// NewBlockchain create a new blockchain with the first block is genesis block
//	if the blockchain already exists, return it
func NewBlockchain() *Blockchain {
	var tip []byte
	// open a boltDB file.
	//  Note: No error will be raised if the file is missing
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	// Operations on DB are in form of transaction
	// two types are available: read-only & read-write
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil { //if the block doesn't exists
			// create genesis block and put it in the db
			//	update the last block hash which is stored with key 1
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				fmt.Println("Create Bucket Error:", err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				fmt.Println("Put Block Error:", err)
			}
			err = b.Put([]byte("1"), genesis.Hash)
			if err != nil {
				fmt.Println("Put Tip Error", err)
			}
			tip = genesis.Hash
		} else { // if the blockchain already exists, get the hash of the last block
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc
}
