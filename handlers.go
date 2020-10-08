package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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
func returnNTransactions(w http.ResponseWriter, r *http.Request, number string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("transactions/"+number, w, r)
	db, connectErr := connect()
	defer db.Close()
	handle("Error creating a DB connection: ", connectErr)
	var graph []byte
	var count int
	_ = db.Get(&count, "SELECT COUNT(*) FROM transactions")
	graph = loadGraphElementsArray(number)
	w.Write([]byte(graph))

}

// returnSingleTransaction This will print a single tx
func returnSingleTransaction(w http.ResponseWriter, r *http.Request, hash string) {
	db, connectErr := connect()
	defer db.Close()
	handle("Error creating a DB connection: ", connectErr)
	var id int
	errCount := db.Get(&id, "SELECT COUNT(*) FROM transactions WHERE tx_hash = $1", hash)
	handle("There was a problem counting the results: ", errCount)
	idVal := 1
	existEval := id == idVal
	if existEval {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		reportRequest("transactions/"+hash, w, r)
		tx := loadSingleTx(hash)
		w.Write(tx)
	} else if !existEval {
		fmt.Printf("cant do it %d", id)
		w.WriteHeader(http.StatusNotFound)
	}
}

func reportRequest(name string, w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Printf(brightgreen+"\n/%s"+white+" by "+brightcyan+"%s\n"+white+"Agent: "+brightcyan+"%s\n"+nc, name, r.RemoteAddr, userAgent)
}

func returnStatsWeb(w http.ResponseWriter, r *http.Request, keys *ED25519Keys) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reportRequest("Stats", w, r)

	db, connectErr := connect()
	defer db.Close()
	handle("Error creating a DB connection: ", connectErr)
	wlPeerCount := countWhitelistPeers()
	graph := []Transaction{}
	_ = db.Select(&graph, "SELECT * FROM transactions")
	numTx := len(graph)
	version := semverInfo()
	pkstring := keys.publicKey

	infoStruct := &StatsDetail{channelName, channelDesc, version, channelCont, pkstring, numTx, wlPeerCount}
	infoJSON, _ := json.Marshal(infoStruct)
	w.Write([]byte(infoJSON))
}

// loadGraphElementsArray outputs the entire graph as an array of Transactions
func loadGraphElementsArray(number string) []byte {
	amount, _ := strconv.Atoi(number)
	if amount <= 300 {
		db, connectErr := connect()
		defer db.Close()
		handle("Error creating a DB connection: ", connectErr)
		graph := []Transaction{}
		err := db.Select(&graph, "SELECT * FROM transactions ORDER BY tx_time DESC LIMIT $1", number)
		graphJSON, _ := json.MarshalIndent(&graph, "", "  ")
		switch {
		case err != nil:
			handle("There was a problem loading the graph: ", err)
			return graphJSON[0:0]
		default:
			return graphJSON
		}
	} else {
		// if they request more than 300, return 300
		number = "300"
		db, connectErr := connect()
		defer db.Close()
		handle("Error creating a DB connection: ", connectErr)
		graph := []Transaction{}
		err := db.Select(&graph, "SELECT * FROM transactions ORDER BY tx_time DESC LIMIT $1", number)
		graphJSON, _ := json.MarshalIndent(&graph, "", "  ")
		switch {
		case err != nil:
			handle("There was a problem loading the graph: ", err)
			return graphJSON[0:0]
		default:
			return graphJSON
		}
	}

}
