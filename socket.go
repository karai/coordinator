package main

import (
	"bytes"
	"crypto/sha512"
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
			break
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
				fmt.Printf("\nCreating cert: %s", trimmedPubKey[:64])
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
						fmt.Printf("\nCert Body: %s", certBody[:192])
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
				trustedSessionParser(conn, keyCollection)
			}

		}
		if bytes.HasPrefix(msg, rtrnMsg) {
			fmt.Printf("\nreturn message: %s", string(msg))

			// strip away the `RTRN` command prefix
			input := strings.TrimLeft(string(msg), "RTRN ")
			trimmedInput := strings.TrimSuffix(input, "\n")
			// fmt.Printf("\ntrimmedInput: %s", trimmedInput)
			var cert = strings.Split(trimmedInput, " ")

			trimmer := strings.TrimSuffix(cert[1], "\n")
			trimmedBytes := []byte(trimmer)

			var hashOfTrimmer = sha512.Sum512(trimmedBytes)
			var encodedHashOfTrimmer = hex.EncodeToString(hashOfTrimmer[:])
			if !verifySignature(encodedHashOfTrimmer, keyCollection.publicKey, cert[0]) {
				fmt.Printf("\nsig doesnt verify")
			}
			if verifySignature(encodedHashOfTrimmer, keyCollection.publicKey, cert[0]) {
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
