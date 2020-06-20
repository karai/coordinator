package main

import (
	"crypto/rand"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
)

type keys struct {
	pubKey    ed25519.PublicKey
	privKey   ed25519.PrivateKey
	signedKey []byte
}

// [✔️] Coord: Generates Secret Key (CA:SK)and Public Key (CA:PK)
// [✔️] Coord: Signs CA:PK with CA:SK(CA:S)
// [❌] Coord: Publishes CA:S & CA:PK in pointer record
// [✔️] Node:  Generates Secret Key (N1:SK) and Public Key(N1:PK)
// [✔️] Node:  Initial Connection Sends N1:PK to Coord
// [✔️] Coord: Signs N1:PK with CA:SK (N1:S)
// [✔️] Coord: Sends CA:N1:S to Node
// [❌] Node:  Verifies N1:S using known CA:PK from pointer (Good Coordinator)
// [ ] Node:  Signs N1:PK with N1:SK (N1:S)
// [ ] Node:  Sends N1:S to Coord
// [ ] Coord: Hashes (N1:S) and signs with CA:SK Node1 Cert (N1:C)
// [ ] Coord: Sends N1:C to Node
// [ ] Node:  Requests Cert Revocation List
// [ ] Coord: Sends CRL to Node

// grab the response from when the client signs on and sends the key over
// take that response and shove it up your ass kek
// verify it usin te well documented golang docs for ed25519 verify
// this will likely need to be done by completingf the connect channel function.
// TODO: make zeus info into variables
// fgor now make thed channel connection function just connect to zeus.

// generateEd25519 Generate credentials and if Coordinator then sign them
// https://hackmd.io/@ZL2uKk4cThC4TG0z7Wu7sg/H1Ubn6d9L
func generateEd25519() *keys {
	keyCollection := keys{}

	// Generate our primary public and private key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)

	keyCollection.pubKey = pubKey
	keyCollection.privKey = privKey
	keyCollection.signedKey = ed25519.Sign(keyCollection.privKey, keyCollection.pubKey)

	handle("Something went wrong generating a key: ", err)

	// Delete stored public and private keys
	deleteFile("pub.key")
	deleteFile("priv.key")

	// Create files read to store new key data
	createFile("pub.key")
	createFile("priv.key")

	logrus.Info("P2P Public Key: ")
	fmt.Printf("%x\n", keyCollection.pubKey)

	logrus.Info("P2P Private Key: ")
	fmt.Printf("%x\n", keyCollection.privKey)

	logrus.Info("P2P Signed Pubkey: ")
	fmt.Printf("%x\n", keyCollection.signedKey)

	// Write new pub / priv key to file
	writeFileBytes("pub.key", keyCollection.pubKey)
	writeFileBytes("priv.key", keyCollection.privKey)

	return &keyCollection
}

/*
func coordSignNodePubKey(nodePubKey ed25519.PublicKey) []byte {
	privKey = readFileBytes("priv.key")
	trimmedPrivKey = privKey[:64]
	fmt.Printf("Coord Private Key: %x\n", trimmedPrivKey)
	fmt.Printf("Node Pub Key: %x\n", []byte(nodePubKey))
	signedNodePubKey := ed25519.Sign(trimmedPrivKey, trimmedPubKey)
	fmt.Printf("P2P Signed Pubkey: %x\n", signedNodePubKey)
	return signedNodePubKey
}
*/
