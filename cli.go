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
	cli.validateArgs()
	createBlockchain := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	createBlockchainData := createBlockchain.String("address", "", "Address of transaction")
	getBalanceData := getBalanceCmd.String("address", "", "address to get balance")
	sendFrom := sendCmd.String("from", "", "from who")
	sendTo := sendCmd.String("to", "", "send to")
	sendAmount := sendCmd.String("amount", "", "Amount to send")

	switch os.Args[1] {
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
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic("Error parsing print chain:", err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic("Error parsing print chain:", err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchain.Parsed() {
		if *createBlockchainData == "" {
			createBlockchain.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainData)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceData)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == "" {
			sendCmd.Usage()
			os.Exit(1)
		}
		amountToSend, _ := strconv.Atoi(*sendAmount)
		cli.send(*sendFrom, *sendTo, amountToSend)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func (cli *CLI) createBlockchain(data string) {
	cli.bc = NewBlockchain(data)

}

func (cli *CLI) printChain() {
	cli.bc = NewBlockchain("")
	defer cli.bc.db.Close()
	bci := cli.bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
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

func (cli *CLI) getBalance(address string) {
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Transaction completed!")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
