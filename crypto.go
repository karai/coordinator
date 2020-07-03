package main

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/ed25519"
)

// ED25519Keys This is a struct for holding keys and a signature.
type ED25519Keys struct {
	publicKey  string
	privateKey string
	signedKey  string
}

func generateKeys() *ED25519Keys {
	keys := ED25519Keys{}

	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		panic(err)
	}

	keys.privateKey = hex.EncodeToString(privKey[0:32])

	keys.publicKey = hex.EncodeToString(pubKey)

	signedKey := ed25519.Sign(privKey, pubKey)

	keys.signedKey = hex.EncodeToString(signedKey)

	return &keys
}

func sign(myKeys *ED25519Keys, msg string) string {
	messageBytes := []byte(msg)

	privateKey, err := hex.DecodeString(myKeys.privateKey)

	if err != nil {
		panic(err)
	}

	publicKey, err := hex.DecodeString(myKeys.publicKey)

	if err != nil {
		panic(err)
	}

	privateKey = append(privateKey, publicKey...)

	signature := ed25519.Sign(privateKey, messageBytes)

	return hex.EncodeToString(signature)
}

func signKey(myKeys *ED25519Keys, publicKey string) string {
	messageBytes, err := hex.DecodeString(publicKey)

	if err != nil {
		panic(err)
	}

	privateKey, err := hex.DecodeString(myKeys.privateKey)

	if err != nil {
		panic(err)
	}

	pubKey, err := hex.DecodeString(myKeys.publicKey)

	if err != nil {
		panic(err)
	}

	privateKey = append(privateKey, pubKey...)

	signature := ed25519.Sign(privateKey, messageBytes)

	return hex.EncodeToString(signature)
}

func verifySignature(publicKey string, msg string, signature string) bool {
	pubKey, err := hex.DecodeString(publicKey)

	if err != nil {
		panic(err)
	}

	messageBytes := []byte(msg)

	sig, err := hex.DecodeString(signature)

	if err != nil {
		panic(err)
	}

	return ed25519.Verify(pubKey, messageBytes, sig)
}

func verifySignedKey(publicKey string, publicSigningKey string, signature string) bool {
	pubKey, err := hex.DecodeString(publicKey)

	if err != nil {
		panic(err)
	}

	pubSignKey, err := hex.DecodeString(publicSigningKey)

	if err != nil {
		panic(err)
	}
	sig, err := hex.DecodeString(signature)

	if err != nil {
		panic(err)
	}

	return ed25519.Verify(pubSignKey, pubKey, sig)
}
