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

	"github.com/fatih/color"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/acme/autocert"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
}

var nodePubKeySignature []byte
var joinMsg []byte = []byte("JOIN")
var castMsg []byte = []byte("CAST")
var peerMsg []byte = []byte("PEER")
var pubkMsg []byte = []byte("PUBK")
var nsigMsg []byte = []byte("NSIG")
var tsxnMsg []byte = []byte("TSXN")

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
	api.HandleFunc("/transaction/send", transactionHandler).Methods(http.MethodPost)
	api.HandleFunc("/channel", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		color.Set(color.FgHiGreen, color.Bold)
		fmt.Printf("\n[%s] [%s] Peer socket opened!\n", timeStamp(), conn.RemoteAddr())
		color.Set(color.FgWhite, color.Bold)
		channelAuthAgent(conn, keyCollection, graph)
	})
	if !wantsHTTPS {
		http.ListenAndServe(":"+strconv.Itoa(karaiAPIPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
	}
	if wantsHTTPS {
		http.Serve(autocert.NewListener(sslDomain), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
	}
}

func channelAuthAgent(conn *websocket.Conn, keyCollection *ED25519Keys, graph *Graph) {
	for {
		defer conn.Close()
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			color.Set(color.FgHiYellow, color.Bold)
			fmt.Printf("\n[%s] [%s] Peer socket closed!\n", timeStamp(), conn.RemoteAddr())
			color.Set(color.FgWhite)
			break
		}
		// defer conn.Close()
		if bytes.HasPrefix(msg, pubkMsg) {
			conn.WriteMessage(msgType, []byte(keyCollection.publicKey))
		}
		if bytes.HasPrefix(msg, peerMsg) {
			conn.WriteMessage(msgType, []byte(getPeerID()))
		}
		if bytes.HasPrefix(msg, joinMsg) {
			trimNewline := bytes.TrimRight(msg, "\n")
			trimmedPubKey := bytes.TrimLeft(trimNewline, "JOIN ")
			if len(trimmedPubKey) == 64 {
				var regValidate bool
				regValidate, _ = regexp.MatchString(`[a-f0-9]{64}`, string(trimmedPubKey))
				if regValidate == false {
					fmt.Printf("\nContains illegal characters")
					conn.Close()
					return
				}
				pubkey := string(trimmedPubKey)
				var whitelistPeerCertFile = p2pConfigDir + "/whitelist/" + pubkey + ".cert"
				var blacklistPeerCertFile = p2pConfigDir + "/blacklist/" + pubkey + ".cert"
				if fileExists(blacklistPeerCertFile) {
					data := []byte("Error 403: You are banned")
					conn.WriteMessage(1, data)
					conn.Close()
					return
				}
				if !fileExists(whitelistPeerCertFile) {
					color.Set(color.FgWhite)

					// Sending Coord Pubkey
					signedNodePubKey := signKey(keyCollection, pubkey)

					// When a peer is banned, this will trigger when they try to reconnect
					_ = conn.WriteMessage(msgType, []byte(signedNodePubKey))

					// Read the request which should be for a pubkey
					_, coordRespPubKey, _ := conn.ReadMessage()
					if bytes.HasPrefix(coordRespPubKey, pubkMsg) {
						conn.WriteMessage(msgType, []byte(keyCollection.publicKey))

						// Wait for a request that should look like
						// NSIG + n1 signature
						_, receivedN1s, _ := conn.ReadMessage()
						if bytes.HasPrefix(receivedN1s, nsigMsg) {
							var stringN1S = string(receivedN1s)
							var trimmedStringRecN1S = strings.TrimRight(stringN1S, "\n")
							var n1sToHash = []byte(trimmedStringRecN1S)
							var hashOfN1s = sha512.Sum512(n1sToHash)
							var encodedN1sHash = hex.EncodeToString(hashOfN1s[:])
							var signedHashOfN1s = sign(keyCollection, encodedN1sHash)
							var certMsg = "CERT " + signedHashOfN1s
							var trimmedCertMsg = strings.TrimLeft(certMsg, " ")
							conn.WriteMessage(msgType, []byte(trimmedCertMsg))

							color.Set(color.FgHiGreen, color.Bold)
							fmt.Printf("[%s] [%s] Certificate Granted!\n", timeStamp(), conn.RemoteAddr())
							color.Set(color.FgHiCyan, color.Bold)
							fmt.Printf("user> ")
							color.Set(color.FgHiBlack, color.Bold)
							fmt.Printf("%s\n", pubkey)
							color.Set(color.FgHiRed, color.Bold)
							fmt.Printf("cert> ")
							color.Set(color.FgHiBlack, color.Bold)
							fmt.Printf("%s\n", signedHashOfN1s)
							color.Set(color.FgWhite)
							// Does a peer file for this node exist?
							var peerCertFile = p2pConfigDir + "/whitelist/" + pubkey + ".cert"
							if !fileExists(peerCertFile) {
								createFile(peerCertFile)
								writeFile(peerCertFile, signedHashOfN1s)
							}
						}

					} else {
						conn.Close()
						return
					}
				} else {
					var wbPeerMsg = "Welcome back " + pubkey[:8]
					conn.WriteMessage(1, []byte(wbPeerMsg))

					// Pass the socket to the incoming message parser and invite our new
					// guest in for some tea
					incMsgParser(conn, keyCollection, graph)
				}
			}
		}
	}
}

func processMsg(msg []byte, graph *Graph) bool {
	if bytes.HasPrefix(msg, tsxnMsg) {
		trimMsg := bytes.TrimRight(msg, "\n")
		dataBytes := bytes.TrimLeft(trimMsg, "TSXN ")
		data := string(dataBytes)
		if validJSON(data) {
			graph.addTx(2, string(data))
			return true
		}
	}
	return false
}

func incMsgParser(conn *websocket.Conn, keyCollection *ED25519Keys, graph *Graph) {
	fmt.Printf("We have reached the tx handler")
	for {
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			color.Set(color.FgHiYellow, color.Bold)
			fmt.Printf("\n[%s] [%s] socket: %s\n", timeStamp(), conn.RemoteAddr(), err)
			color.Set(color.FgWhite)
			break
		}
		// processMsg(msg)
		if processMsg(msg, graph) {
			color.Set(color.FgHiGreen, color.Bold)
			fmt.Printf("\n[%s] [%s] Tx Good: %s\n", timeStamp(), conn.RemoteAddr(), err)
			color.Set(color.FgWhite)
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

// transactionHandler When a transaction is sent from a client,
// it goes to the CC first. The CC should then triage and
// validate that transaction, timestamp it and append to a subgraph
func transactionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: triage transactions
	// this should work hand in hand with subgraphConstructor
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
	// peerFile, err := os.OpenFile(configPeerIDFile,
	// 	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// handle("Can't find peer.id file: ", err)
	// defer peerFile.Close()

	// fileToRead, err := ioutil.ReadFile(configPeerIDFile)
	// // fmt.Println(fileToRead)
	// handle("Error: ", err)
	// w.Write([]byte("{\"p2p_peer_ID\": \"" + string(fileToRead) + "\"}"))
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
