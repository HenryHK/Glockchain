package main

import (
	"bytes"
)

// TxInput defines the structure of a transaction input
type TxInput struct {
	Txid      []byte // refers to the transaction the input consumed
	Vout      int    // refers to the index of comsumed outputs within the transaction
	Signature []byte // Sig of this input
	PubKey    []byte // PubKey of creator
}

// UsesKey checks whether the address/public key initiated the transaction
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
