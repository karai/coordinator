package main

import (
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
