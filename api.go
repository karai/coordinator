package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// restAPI() This is the main API that is activated when isCoord == true
func restAPI(keyCollection *ED25519Keys, graph *Graph) {

	// CORS
	corsAllowedHeaders := []string{
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Cache-Control",
		"Content-Security-Policy",
		"Feature-Policy",
		"Referrer-Policy",
		"X-Requested-With"}

	corsOrigins := []string{
		"*",
		"127.0.0.1"}

	corsMethods := []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"OPTIONS"}

	headersCORS := handlers.AllowedHeaders(corsAllowedHeaders)
	originsCORS := handlers.AllowedOrigins(corsOrigins)
	methodsCORS := handlers.AllowedMethods(corsMethods)

	// Init API
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// Home
	api.HandleFunc("/", home).Methods(http.MethodGet)

	// Version
	api.HandleFunc("/version", returnVersion).Methods(http.MethodGet)

	// Peer
	api.HandleFunc("/peer", func(w http.ResponseWriter, r *http.Request) {
		returnPeerID(w, r, keyCollection.publicKey)
	}).Methods(http.MethodGet)

	// Stats
	api.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		returnStats(w, r, keyCollection, graph)
	}).Methods(http.MethodGet)

	// Transactions
	api.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		returnTransactions(w, r, graph)
	}).Methods(http.MethodGet)

	// Transaction by ID
	api.HandleFunc("/transaction/{txid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		txid := vars["txid"]
		reportRequest("transaction/"+txid, w, r)
		returnSingleTransaction(w, r, graph, txid)
	}).Methods(http.MethodGet)

	// Channel Socket
	api.HandleFunc("/channel", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		fmt.Printf(brightgreen+"\n[%s] [%s] Peer socket opened!\n"+white, timeStamp(), conn.RemoteAddr())
		socketAuthAgent(conn, keyCollection, graph)
	})

	// Serve via HTTP
	http.ListenAndServe(":"+strconv.Itoa(karaiAPIPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

// StatsDetail is an object containing strings relevant to the status of a coordinator node.
type StatsDetail struct {
	ChannelName        string `json: "stats_channel_name"`
	ChannelDescription string `json: "stats_channel_description"`
	Version            string `json: "stats_karai_version"`
	ChannelContact     string `json: "stats_channel_contact"`
	PubKeyString       string `json: "stats_pubkey"`
	TxObjectsOnDisk    string `json: "stats_tx_on_disk"`
	TxObjectsInMemory  string `json: "stats_tx_in_memory"`
}

func returnStats(w http.ResponseWriter, r *http.Request, keys *ED25519Keys, graph *Graph) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("Stats", w, r)
	version := semverInfo()
	pkstring := keys.publicKey
	txObjectsOnDisk := countFilesOnDisk(graphDir)
	txObjectsInMemory := countFilesInMemory(graph)
	infoStruct := &StatsDetail{channelName, channelDescription, version, channelContact, pkstring, txObjectsOnDisk, txObjectsInMemory}
	infoJSON, _ := json.Marshal(infoStruct)
	w.Write([]byte(infoJSON))
}
