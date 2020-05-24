package main

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/ed25519"
)

var pubKey []byte
var privKey []byte
var signedPubKey []byte

func generateEd25519() {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	handle("Something went wrong generating a key: ", err)
	trimmedPrivate := privKey[:32]
	fmt.Printf("\npubkey: %x\n", pubKey)
	fmt.Printf("\nprivkey: %x\n", privKey)
	fmt.Printf("\ntprivkey: %x\n", trimmedPrivate)

	signedMsg := ed25519.Sign(privKey, pubKey)

	fmt.Printf("\nmsg: %x\n", signedMsg)
}
