package main

import (
	"crypto/rand"
	"fmt"

	"github.com/sirupsen/logrus"
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
	shortPrivKey := privKey[:32]
	logrus.Info("P2P Public Key: ")
	fmt.Printf("%x\n", pubKey)
	logrus.Info("P2P Private Key: ")
	fmt.Printf("%x\n", privKey)
	logrus.Info("P2P Short Private Key: ")
	fmt.Printf("%x\n", shortPrivKey)
	if isCoordinator == true {
		signedMsg := ed25519.Sign(privKey, pubKey)
		logrus.Info("P2P Signed Message: ")
		fmt.Printf("%x\n", signedMsg)
	}
}
