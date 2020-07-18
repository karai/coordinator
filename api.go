package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// restAPI() This is the main API that is activated when isCoord == true
func restAPI(keyCollection *ED25519Keys, graph *Graph) {
	headersCORS := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin", "Cache-Control", "Content-Security-Policy", "Feature-Policy", "Referrer-Policy", "X-Requested-With"})
	originsCORS := handlers.AllowedOrigins([]string{
		"*",
		"127.0.0.1"})
	methodsCORS := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", home).Methods(http.MethodGet)
	api.HandleFunc("/peer", returnPeerID).Methods(http.MethodGet)
	api.HandleFunc("/version", returnVersion).Methods(http.MethodGet)
	api.HandleFunc("/transactions", returnTransactions).Methods(http.MethodGet)
	api.HandleFunc("/channel", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		fmt.Printf(brightgreen+"\n[%s] [%s] Peer socket opened!\n"+white, timeStamp(), conn.RemoteAddr())
		socketAuthAgent(conn, keyCollection, graph)
	})
	// if !wantsHTTPS {
		http.ListenAndServe(":"+strconv.Itoa(karaiAPIPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
	// }
	// if wantsHTTPS {
	// 	http.Serve(autocert.NewListener(sslDomain), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
	// }
}

func socketAuthAgent(conn *websocket.Conn, keyCollection *ED25519Keys, graph *Graph) {
	for {
		defer conn.Close()
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf(brightyellow+"\n[%s] [%s] Peer socket closed!\n"+white, timeStamp(), conn.RemoteAddr())
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
			fmt.Println("return message: ", string(msg))

			// strip away the `RTRN` command prefix
			input := strings.TrimLeft(string(msg), "RTRN ")
			trimmedInput := strings.TrimSuffix(input, "\n")
			fmt.Println("trimmedInput: ", trimmedInput)
			var cert = strings.Split(trimmedInput, " ")

			trimmer := strings.TrimSuffix(cert[1], "\n")
			trimmedBytes := []byte(trimmer)

			var hashOfTrimmer = sha512.Sum512(trimmedBytes)
			var encodedHashOfTrimmer = hex.EncodeToString(hashOfTrimmer[:])
			if !verifySignature(encodedHashOfTrimmer, keyCollection.publicKey, cert[0]) {
				fmt.Println("sig doesnt verify")
			}
			if verifySignature(encodedHashOfTrimmer, keyCollection.publicKey, cert[0]) {
				fmt.Println("sig verifies")
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

func txParser(msg []byte, graph *Graph) bool {
	if bytes.HasPrefix(msg, tsxnMsg) {
		fmt.Printf("transaction found: %s", string(msg))
		trimMsg := bytes.TrimRight(msg, "\n")
		dataBytes := bytes.TrimLeft(trimMsg, "SEND ")
		data := string(dataBytes)
		if zValidJSON(data) {
			graph.addTx(2, string(data))
			return true
		}
		fmt.Printf("MSG Error: %s", string(msg))
	}
	fmt.Printf("MSG Error: %s", string(msg))
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
			fmt.Printf(brightgreen+"\n[%s] [%s] Tx Good: %s\n"+white, timeStamp(), conn.RemoteAddr(), err)
		} else {
			fmt.Printf("\n Oh no, something has gone very wrong..")
		}
	}
}

// initAPI Check if we are running as a coordinator, if we are, start the API
func initAPI(keyCollection *ED25519Keys) {
	if !isCoordinator {
		return
	}
	go restAPI(keyCollection, spawnGraph())
}

// home This is the home route, it can be used as a
// health check to see if a coordinator is responding.
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello " + appName))
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
func returnPeerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var peerID = getPeerID()
	w.Write([]byte("{\"p2p_peer_ID\": \"" + peerID + "\"}"))
}

// returnVersion This is a dedicated endpoint for returning
// the version as a JSON object
func returnVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"karai_version\": \"" + semverInfo() + "\"}"))
}

// returnTransactions This will print the contents of all of
// the trasnsactions in the graph as an array of JSON objects.
// The {} at the end was a hack because, in a hurry, I
// manually constructed the JSON objects and never went back
// to write proper object creation.
func returnTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	matches, _ := filepath.Glob(graphDir + "/*.json")
	w.Write([]byte("[\n\t"))
	for _, match := range matches {
		w.Write([]byte(printTx(match)))
	}
	w.Write([]byte("{}"))
	w.Write([]byte("\n]"))
}

func getPeerID() string {
	peerFile, err := os.OpenFile(configPeerIDFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handle("Can't find peer.id file: ", err)
	defer peerFile.Close()
	fileToRead, err := ioutil.ReadFile(configPeerIDFile)
	var peerID = string(fileToRead)
	fmt.Println(peerID)
	handle("Error: ", err)
	return peerID
}
