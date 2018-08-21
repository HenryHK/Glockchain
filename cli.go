package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI defines the structure of CLI interface
type CLI struct {
	bc *Blockchain
}

// Run simply runs CLI struct
func (cli *CLI) Run() {
	// cli.validateArgs()
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	createBlockchain := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")
	createBlockchainData := createBlockchain.String("address", "", "Address of transaction")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic("Error parsing add block:", err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic("Error parsing print chain:", err)
		}
	case "createblockchain":
		err := createBlockchain.Parse(os.Args[2:])
		if err != nil {
			log.Panic("Error parsing print chain:", err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if createBlockchain.Parsed() {
		if *createBlockchainData == "" {
			createBlockchain.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainData)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func (cli *CLI) addBlock(data string) {
	cli.bc.MineBlock(data)
	fmt.Printf("Successfully add block!")
}

func (cli *CLI) createBlockchain(data string) {
	bc := NewBlockchain(data)
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.HashTransactions())
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Add Block to Blockchain: Glockchain addblock [DATA]")
	fmt.Println("Print blockchain: Glockchain printchain")
}
