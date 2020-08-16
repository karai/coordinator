package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

func socketAuthAgent(conn *websocket.Conn, keyCollection *ED25519Keys, graph *Graph) {
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
				trustedSessionParser(conn, keyCollection, graph)
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

func addTransactions(number int, graph *Graph) {
	sum := 0
	for i := 1; i < number; i++ {
		sum += i
		msg := string(unixTimeStampNano())
		graph.addTx(2, msg)
	}
}

func txParser(msg []byte, graph *Graph) bool {
	trimMsg := bytes.TrimRight(msg, "\n")
	data := string(trimMsg)
	if validJSON(data) {
		fmt.Printf("\nSubmitting transaction: %s", data)
		graph.addTx(2, string(data))
		return true
	}
	fmt.Printf("\nJSON Error: %s", string(msg))
	return false
}

func trustedSessionParser(conn *websocket.Conn, keyCollection *ED25519Keys, graph *Graph) {
	fmt.Printf(brightgreen+"\n[%s] [%s] Trusted Socket Session OPEN"+white, timeStamp(), conn.RemoteAddr())
	for {
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf(brightyellow+"\n[%s] [%s] socket: %s\n"+white, timeStamp(), conn.RemoteAddr(), err)
			break
		}
		if txParser(msg, graph) {
			fmt.Printf(brightgreen+"\n[%s] [%s] Tx Good! \n"+white, timeStamp(), conn.RemoteAddr())

		} else {
			fmt.Printf("\n Oh no, something has gone very wrong..\n %s", msg)
		}
	}
}

// initAPI Check if we are running as a coordinator, if we are, start the API
func initAPI(keyCollection *ED25519Keys) *Graph {
	graph := spawnGraph()
	go restAPI(keyCollection, graph)
	return graph
}

// home This is the home route, it can be used as a
// health check to see if a coordinator is responding.
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("", w, r)
	w.Write([]byte("\"" + appName + "\""))
}

// notfound when an API route is unrecognized, we should reply with
// something to communicate that.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"bruh": "lol"}`))
}

// returnPeerID will deliver the contents of peer.id file
// through the API. This is the first step in connecting to
// a tx channel.
func returnPeerID(w http.ResponseWriter, r *http.Request, pubkey string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("peer", w, r)
	// var peerID = readFile(pubKeyFilePath)
	w.Write([]byte("{\"p2p_peer_ID\": \"" + pubkey + "\"}"))
}

// returnVersion This is a dedicated endpoint for returning
// the version as a JSON object
func returnVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("version", w, r)
	w.Write([]byte("{\"karai_version\": \"" + semverInfo() + "\"}"))
}

// returnTransactions This will print the contents of all of
// the trasnsactions in the graph as an array of JSON objects.
func returnTransactions(w http.ResponseWriter, r *http.Request, graph *Graph) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("transactions", w, r)
	graphJSON := prettyPrintGraphJSON(graph)
	w.Write([]byte(graphJSON))

}

func reportRequest(name string, w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Printf(brightgreen+"\n/%s"+white+" by "+brightcyan+"%s\n"+white+"Agent: "+brightcyan+"%s\n"+nc, name, r.RemoteAddr, userAgent)
}

// returnTransactions This will print the contents of all of
// the trasnsactions in the graph as an array of JSON objects.
func returnSingleTransaction(w http.ResponseWriter, r *http.Request, graph *Graph, transaction string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("transaction/"+transaction, w, r)
	singleTransaction := printTx(graphDir + "/" + transaction + ".json")
	w.Write([]byte(singleTransaction))
}
