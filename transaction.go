package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10

// Transaction defines the structure of a transaction in our blockchain
type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// TxOutput defines the structure of a transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// TxInput defines the structure of a transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// SetID sets ID of a transaction, it's a hash of a transaction itself
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx Transaction) isCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TxInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX creates new coinbase transaction and return its pointer
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{subsidy, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()
	return &tx
}

// NewUTXOTransaction creates new utxo transaction
// func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
// 	var inputs []TxInput
// 	var outputs []TxOutput

// 	acc, validOutputs := bc.
// }
