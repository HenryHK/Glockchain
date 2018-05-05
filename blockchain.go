package main

// Blockchain is the chain holding blocks
type Blockchain struct {
	blocks []*Block
}

// AddBlock is to add a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewBlockchain create a new blockchain with the first block is genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
