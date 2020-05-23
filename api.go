package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

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
	// logrus.Error(http.ListenAndServe(":"+strconv.Itoa(karaiPort), r))
	// logrus.Error(http.ListenAndServe(":"+strconv.Itoa(karaiPort), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api)))
	logrus.Error(http.Serve(autocert.NewListener(sslDomain), handlers.CORS(headersCORS, originsCORS, methodsCORS)(api)))
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

	peerFile, err := os.OpenFile(configPeerIDFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handle("Can't find peer.id file: ", err)
	defer peerFile.Close()

	fileToRead, err := ioutil.ReadFile(configPeerIDFile)
	// fmt.Println(fileToRead)
	handle("Error: ", err)
	w.Write([]byte("{\"p2p_peer_ID\": \"" + string(fileToRead) + "\"}"))
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
