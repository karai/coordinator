package main

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/ed25519"
)

var pubKey []byte
var privKey []byte
var signedPubKey []byte

// generateEd25519 Generate credentials and if Coordinator then sign them
// https://hackmd.io/@ZL2uKk4cThC4TG0z7Wu7sg/H1Ubn6d9L
func generateEd25519() {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	handle("Something went wrong generating a key: ", err)
	trimmedPrivate := privKey[:32]
	if isCoordinator == true {
		signedMsg := ed25519.Sign(privKey, pubKey)
		fmt.Printf("\nmsg: %x\n", signedMsg)
	}
	fmt.Printf("\npubkey: %x\n", pubKey)
	fmt.Printf("\nprivkey: %x\n", privKey)
	fmt.Printf("\ntprivkey: %x\n", trimmedPrivate)
}
