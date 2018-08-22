package main

// TxOutput defines the structure of a transaction output
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock simply locks an output
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
