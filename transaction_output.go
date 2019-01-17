package main

import (
	"bytes"
)

// TxOutput defines the structure of a transaction output
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock simply locks an output, using PubKey
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	// remove version and checksum
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey chekcs if provided public key hash was used to lock the output
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTxOutput create a TxOuput
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}
