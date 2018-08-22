package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressCheckSumLen = 4

// Wallet is a key pair identify a specific user in blockchain
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates new wallet (new key pair actually)
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic("Error generating ecdsa key")
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

// GetAddress returns a wallet's address. The address is derived from public key
func (w Wallet) GetAddress() []byte {
	// double hashing the pub key
	pubKeyHash := HashPubKey(w.PublicKey)

	// prepend version
	versionedPayload := append([]byte{version}, pubKeyHash...)

	// calculate the checksum
	checksum := checksum(versionedPayload)

	// combine all together
	fullPayload := append(versionedPayload, checksum...)

	// use base58 encoded
	address := Base58Encode(fullPayload)

	return address
}

// HashPubKey returns a double-hashed public key
func HashPubKey(pubKey []byte) []byte {
	// hash public key for the first time using SHA256
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	// hash pubic key for the second time using RIPEMD160
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressCheckSumLen]
}
