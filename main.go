package main

import (
	"fmt"
)

func main() {
	blockchain := NewBlockchain()

	blockchain.AddBlock("The first added block.")
	blockchain.AddBlock("The second added block.")

	for _, block := range blockchain.blocks {
		fmt.Printf("Prev. hash %x\n", block.PrevBlockHash)
		fmt.Printf("Data %s\n", block.Data)
		fmt.Printf("Hash %x\n", block.Hash)
		fmt.Println("--------")
	}
}
