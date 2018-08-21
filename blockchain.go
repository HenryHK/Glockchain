package main

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

const genesisCoinbaseData = "Make Australian Great Again"
const blocksBucket = "blocks"
const dbFile = "blockchain.db"

// Blockchain is the chain holding blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// BlockchainIterator stores the information to iterate persistent DB
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Iterator create BlockchainIterator from Blockchain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
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

// MineBlock is to add a new block to the blockchain
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic("Error get last hash from db:", err)
	}

	newBlock := NewBlock(transactions, lastHash)

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

// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	// make a list of unspent transactions
	var unspentTxs []Transaction
	// make a map to store spent transactions
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		// traverse each transaction in a block
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			// traverse the outputs in a transaction
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil { // check if this output is spent
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx { // the output is spent
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.isCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTxs
}

// NewBlockchain create a new blockchain with the first block is genesis block. if there exists blockchain already,do nothing.
func NewBlockchain(address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic("Error opening db file:", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)
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
