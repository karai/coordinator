package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

func socketAuthAgent(conn *websocket.Conn, keyCollection *ED25519Keys) {
	for {
		defer conn.Close()
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer disconnected\n"+white, timeStamp(), conn.RemoteAddr())
			return
		}
		if bytes.HasPrefix(msg, joinMsg) {
			msgToString := string(msg)
			trimNewlinePrefix := strings.TrimRight(msgToString, "\n")
			trimmedPubKey := strings.TrimLeft(trimNewlinePrefix, "JOIN ")
			var regValidate bool
			regValidate, _ = regexp.MatchString(`[a-f0-9]{64}`, trimmedPubKey[:64])
			if regValidate == false {
				fmt.Printf("\nContains illegal characters")
				conn.Close()
				return
			}
			var whitelistPeerCertFile = p2pConfigDir + "/whitelist/" + trimmedPubKey + ".cert"
			var blacklistPeerCertFile = p2pConfigDir + "/blacklist/" + trimmedPubKey + ".cert"
			if fileExists(blacklistPeerCertFile) {
				data := []byte("Error 403: You are banned")
				conn.WriteMessage(1, data)
				conn.Close()
				return
			}
			if !fileExists(whitelistPeerCertFile) {
				fmt.Printf("\nCreating cert: %s.cert", trimmedPubKey[:64])
				capkMsgString := string(capkMsg) + " " + keyCollection.publicKey
				_ = conn.WriteMessage(msgType, []byte(capkMsgString))
				_, receiveNCAS, _ := conn.ReadMessage()
				receiveNCASString := string(receiveNCAS)
				if strings.HasPrefix(receiveNCASString, string(ncasMsg)) {
					convertString := string(receiveNCAS)
					trimNewLine := strings.TrimRight(convertString, "\n")
					nCASig := strings.TrimPrefix(trimNewLine, "NCAS ")
					if verifySignedKey(keyCollection.publicKey[:64], trimmedPubKey[:64], nCASig) {
						signNpkWithCask := signKey(keyCollection, nCASig[:64])
						certForN := string(certMsg) + " " + nCASig[:64] + signNpkWithCask[:128]
						certBody := nCASig[:64] + signNpkWithCask[:128]
						_ = conn.WriteMessage(1, []byte(certForN))
						createFile(whitelistPeerCertFile)
						writeFile(whitelistPeerCertFile, certBody[:192])
						fmt.Printf("\nCert Name: %s", whitelistPeerCertFile)
						fmt.Printf("\nCert Body: %s\n", certBody[:192])
						sessionAgent(conn, trimmedPubKey)
					}
					if !verifySignedKey(keyCollection.publicKey[:64], trimmedPubKey[:64], nCASig) {
						fmt.Printf("\nSignature does not verify!")
						conn.Close()
					}

				}
				// Read the request which should be for a pubkey
				_, coordRespPubKey, _ := conn.ReadMessage()
				if bytes.HasPrefix(coordRespPubKey, pubkMsg) {
					conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
				}
			} else {
				var wbPeerMsg = "WCBK " + trimmedPubKey[:8]
				conn.WriteMessage(1, []byte(wbPeerMsg))
				sessionAgent(conn, trimmedPubKey)
			}

		}
		if bytes.HasPrefix(msg, rtrnMsg) {
			fmt.Printf("\nreturn message: %s", string(msg))

			input := strings.TrimLeft(string(msg), "RTRN ")
			var cert = strings.Split(input, " ")

			if !verifySignature(cert[0], keyCollection.publicKey, cert[1]) {
				fmt.Printf("\nsig doesnt verify")
			}
			if verifySignature(cert[0], keyCollection.publicKey, cert[1]) {
				fmt.Printf("\nsig verifies")
			}
		}
		if bytes.HasPrefix(msg, pubkMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		}
		if bytes.HasPrefix(msg, peerMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		}

	}
}

func sessionAgent(conn *websocket.Conn, sessionPubKey string) {
	fmt.Printf(brightgreen+"\n[%s] [%s] New socket session"+white, timeStamp(), conn.RemoteAddr())
	for {
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf(brightyellow+"\n[%s] [%s] socket: %s\n"+white, timeStamp(), conn.RemoteAddr(), err)
			break
		}
		if txParser(msg) {
			fmt.Printf(brightgreen+"\n[%s] [%s] Tx Good! \n"+white, timeStamp(), conn.RemoteAddr())
			okbyte := []byte("OK " + sessionPubKey)
			_ = conn.WriteMessage(1, okbyte)
		} else if !txParser(msg) {
			fmt.Printf("\nThis transaction has failed to parse:\n %s", msg)
		}
	}
}

func txParser(msg []byte) bool {
	// SEND
	if bytes.HasPrefix(msg, sendMsg) {

		// Remove the SEND command
		input := strings.TrimLeft(string(msg), "SEND ")

		// Split the remaining texts
		var sendString = strings.Split(input, " ")

		// Use the sender pubkey to look for a matching cert
		pubkey := sendString[0]
		fileNamePubKey := pubkey[:64]
		fileName := p2pWhitelistDir + "/" + fileNamePubKey + ".cert"
		readCert := readFile(fileName)

		// If the cert in our records matches the cert the sender gave us,
		// the sender is verified.
		certChecksOut := readCert == sendString[1]
		if !certChecksOut {
			fmt.Printf("\nTx signature does not verify.")
			return false
		}
		if certChecksOut {
			decodedBytes, _ := hex.DecodeString(sendString[2])
			transactionBody := strings.TrimRight(string(decodedBytes), "\n")
			fmt.Printf(brightgreen+"\nTX signature verified for pubkey %s\n%s", brightcyan+pubkey[:4]+"..."+pubkey[60:64]+nc, brightwhite+transactionBody+nc)
			createTransaction("2", transactionBody)
		}
	}
	return true
}
