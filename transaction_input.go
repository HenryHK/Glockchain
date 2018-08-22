package main

// TxInput defines the structure of a transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// UsesKey checks thaty an input usese a specific key to unlock an output
// Notice: inputs store raw public keys (not hashed)
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
