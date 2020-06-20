package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/fatih/color"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/ed25519"
)

var upgrader = websocket.Upgrader{
	// EnableCompression: true,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var joinMsg []byte = []byte("JOIN")
var castMsg []byte = []byte("CAST")
var versionMsg []byte = []byte("VERSION")
var peerMsg []byte = []byte("PEER")

// restAPI() This is the main API that is activated when isCoord == true
func restAPI() {
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
		channelSocketHandler(conn)
	})
	if !wantsHTTPS {
		logrus.Debug(http.ListenAndServe(":"+strconv.Itoa(karaiAPIPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api)))
	}
	if wantsHTTPS {
		logrus.Debug(http.Serve(autocert.NewListener(sslDomain), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api)))
	}
}

func channelSocketHandler(conn *websocket.Conn) {

	// TODO: look into whether it makes sense to use channels for concurrency
	// in any of this.

	for {
		defer conn.Close()
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			color.Set(color.FgHiYellow, color.Bold)
			fmt.Printf("\n[%s] [%s] Peer socket closed!", timeStamp(), conn.RemoteAddr())
			color.Set(color.FgHiWhite, color.Bold)
			break
		}
		defer conn.Close()
		if bytes.HasPrefix(msg, joinMsg) {
			trimNewline := bytes.TrimRight(msg, "\n")
			trimmedPubKey := bytes.TrimLeft(trimNewline, "JOIN ")
			if len(trimmedPubKey) == 64 {
				var regValidate bool
				regValidate, _ = regexp.MatchString(`[a-f0-9]{64}`, string(trimmedPubKey))
				if regValidate == false {
					logrus.Error("Contains illegal characters")
					conn.Close()
					return
				}
				fmt.Printf("\n- Node Pub Key Received: %v\n", string(trimmedPubKey))
				privKey = readFileBytes("priv.key")
				trimmedPrivKey = privKey[:64]
				fmt.Printf("- Coord Private Key: %x\n", string(trimmedPrivKey))
				fmt.Printf("- Node Pub Key: %v\n", string(trimmedPubKey))
				signedNodePubKey := ed25519.Sign(trimmedPrivKey, trimmedPubKey)
				fmt.Printf("- P2P Signed Pubkey: %x\n", string(signedNodePubKey))
				hexSig := fmt.Sprintf("%x", signedNodePubKey)
				err = conn.WriteMessage(msgType, []byte(hexSig))
				handle("respond with signed node pubkey", err)

				if !fileExists(p2pConfigDir + "/" + string(trimmedPubKey) + ".pubkey") {
					createFile(p2pConfigDir + "/" + string(trimmedPubKey) + ".pubkey")
					writeFileBytes(p2pConfigDir+"/"+string(trimmedPubKey)+".pubkey", signedNodePubKey)
				}
			} else {
				fmt.Printf("Join PubKey %s has incorrect length. PubKey received has a length of %v", string(trimmedPubKey), len(trimmedPubKey))
				conn.Close()
				return
			}
		}
		if bytes.HasPrefix(msg, joinMsg) {
			trimNewline := bytes.TrimRight(msg, "\n")
			trimmedPubKey := bytes.TrimLeft(trimNewline, "JOIN ")
			if len(trimmedPubKey) == 64 {
				var regValidate bool
				regValidate, _ = regexp.MatchString(`[a-f0-9]{64}`, string(trimmedPubKey))
				if regValidate == false {
					logrus.Error("Contains illegal characters")
					conn.Close()
					return
				}
				fmt.Printf("\n- Node Pub Key Received: %v\n", string(trimmedPubKey))
				privKey = readFileBytes("priv.key")
				trimmedPrivKey = privKey[:64]
				fmt.Printf("- Coord Private Key: %x\n", string(trimmedPrivKey))
				fmt.Printf("- Node Pub Key: %v\n", string(trimmedPubKey))
				signedNodePubKey := ed25519.Sign(trimmedPrivKey, trimmedPubKey)
				fmt.Printf("- P2P Signed Pubkey: %x\n", string(signedNodePubKey))
				hexSig := fmt.Sprintf("%x", signedNodePubKey)
				err = conn.WriteMessage(msgType, []byte(hexSig))
				handle("respond with signed node pubkey", err)

				if !fileExists(p2pConfigDir + "/" + string(trimmedPubKey) + ".pubkey") {
					createFile(p2pConfigDir + "/" + string(trimmedPubKey) + ".pubkey")
					writeFileBytes(p2pConfigDir+"/"+string(trimmedPubKey)+".pubkey", signedNodePubKey)
				}
			} else {
				fmt.Printf("Join PubKey %s has incorrect length. PubKey received has a length of %v", string(trimmedPubKey), len(trimmedPubKey))
				conn.Close()
				return
			}
		}
		if bytes.HasPrefix(msg, versionMsg) {
			conn.WriteMessage(msgType, []byte(semverInfo()))
		}
		if bytes.HasPrefix(msg, peerMsg) {
			conn.WriteMessage(msgType, []byte(getPeerID()))
		}
	}
}

// initAPI Check if we are running as a coordinator, if we are, start the API
func initAPI() {
	if !isCoordinator {
		logrus.Debug("isCoordinator == false, skipping webserver deployment")
	} else {
		go restAPI()
	}
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

// returns the contents of the peer.id file
func getPeerID() string {
	peerFile, err := os.OpenFile(configPeerIDFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handle("Can't find peer.id file: ", err)
	defer peerFile.Close()

	fileToRead, err := ioutil.ReadFile(configPeerIDFile)

	var peerId = string(fileToRead)
	fmt.Println(peerId)

	handle("Error: ", err)

	return peerId
}

// returnPeerID will deliver the contents of peer.id file
// through the API. This is the first step in connecting to
// a tx channel.
func returnPeerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var peerId = getPeerID()

	w.Write([]byte("{\"p2p_peer_ID\": \"" + peerId + "\"}"))
}

// returnVersion This is a dedicated endpoint for returning
// the version as a JSON object
func returnVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"karai_version\": \"" + semverInfo() + "\"}"))
}

func getTransactions() {
	matches, _ := filepath.Glob(graphDir + "/*.json")

	for _, match := range matches {
		fmt.Println(match)
	}
}

// returnTransactions This will print the contents of all of
// the trasnsactions in the graph as an array of JSON objects.
// The {} at the end was a hack because, in a hurry, I
// manually constructed the JSON objects and never went back
// to write proper object creation.
func returnTransactions(w http.ResponseWriter, r *http.Request) {
	getTransactions()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	matches, _ := filepath.Glob(graphDir + "/*.json")
	w.Write([]byte("[\n\t"))
	for _, match := range matches {
		fmt.Printf(printTx(match))
		w.Write([]byte(printTx(match)))
	}
	w.Write([]byte("{}"))
	w.Write([]byte("\n]"))
}
