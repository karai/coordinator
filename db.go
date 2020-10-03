package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/sha3"
)

// Graph is a collection of transactions
type Graph struct {
	Transactions []Transaction `json:"transactions"`
}

// Transaction This is the structure of the transaction
type Transaction struct {
	Time string `json:"time" db:"tx_time"`
	Type string `json:"type" db:"tx_type"`
	Hash string `json:"hash" db:"tx_hash"`
	Data string `json:"data" db:"tx_data"`
	Prev string `json:"prev" db:"tx_prev"`
	Epoc string `json:"epoc" db:"tx_epoc"`
	Subg string `json:"subg" db:"tx_subg"`
	Prnt string `json:"prnt" db:"tx_prnt"`
	Mile bool   `json:"mile" db:"tx_mile"`
	Lead bool   `json:"lead" db:"tx_lead"`
}

// MilestoneData is a struct that defines the contents of the data field to be parsed
type MilestoneData struct {
	TxInterval int
}

// SubGraph defines participants and child transactions of a subgraph
type SubGraph struct {
	PubKeys  []string
	Children []string
}

// connect will create an active DB connection
func connect() (*sqlx.DB, error) {
	connectParams := fmt.Sprintf("user=%s dbname=%s sslmode=%s", dbUser, dbName, dbSSL)
	db, err := sqlx.Connect("postgres", connectParams)
	return db, err
}

// verifyGraph will check each transaction in the graph
func verifyGraph() (bool, error) {
	return true, nil
}

// loadGraphArray outputs the entire graph as an array of Transactions
func loadGraphArray() []byte {
	db, connectErr := connect()
	defer db.Close()
	handle("Error creating a DB connection: ", connectErr)

	graph := []Transaction{}
	err := db.Select(&graph, "SELECT * FROM transactions")
	graphJSON, _ := json.MarshalIndent(&graph, "", "  ")
	switch {
	case err != nil:
		handle("There was a problem loading the graph: ", err)
		return graphJSON
	default:
		return graphJSON
	}
}

// loadGraphElementsArray outputs the entire graph as an array of Transactions
func loadGraphElementsArray(number string) []byte {
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

// addBulkTransactions adds $number of tx to the graph immediately
func addBulkTransactions(number int) {
	sum := 0
	for i := 1; i < number; i++ {
		sum += i
		msg := string(unixTimeStampNano())
		createTransaction("2", msg)
	}
}

// loadSingleTx outputs a single tx by hash
func loadSingleTx(hash string) []byte {
	db, connectErr := connect()
	defer db.Close()
	handle("Error creating a DB connection: ", connectErr)

	tx := Transaction{}
	err := db.Get(&tx, "SELECT * FROM transactions WHERE tx_hash=$1", hash)
	// if &tx.Data == nil {
	txJSON, _ := json.MarshalIndent(&tx, "", "  ")
	switch {
	case err != nil:
		handle("There was a problem loading the graph: ", err)
		return txJSON
	default:
		fmt.Printf(brightcyan+"\nReturning tx with hash:"+brightgreen+"\n%s\n"+nc, hash)
		return txJSON
		// }
	}
}

func generateRandomTransactions() {
	min := 1
	max := 100
	rand.Seed(time.Now().UTC().UnixNano())
	number := rand.Intn(max-min) + min
	fmt.Printf("+%v tx in subgraph %s..  ", number, thisSubgraphShortName)
	addBulkTransactions(number)
	// fmt.Printf(brightgreen+" Done! "+brightyellow+"Sleeping %v seconds."+nc, number)
	if number%3 == 0 {
		time.Sleep(time.Duration(number) * time.Minute)
	}
	time.Sleep(time.Duration(number) * time.Second)
	generateRandomTransactions()
}

func createBenchmark(number int) {
	// start := time.Now()

	addBulkTransactions(10)
	// t := time.Now()
	// elapsed := t.Sub(start)
	// fmt.Printf(brightcyan+"\nBenchmark\nFinished %v operations in %s\n"+nc, number, elapsed)
}

// createRoot Transaction channels start with a rootTx transaction always
func createRoot() error {
	db, connectErr := connect()
	defer db.Close()
	handle("Error creating a DB connection: ", connectErr)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM transactions ORDER BY $1 DESC", "tx_time").Scan(&count)
	switch {
	case err != nil:
		handle("There was a problem counting database transactions: ", err)
		return err
	default:
		//  fmt.Printf("Found %v transactions in the db.", count)
		if count == 0 {
			txTime := unixTimeStampNano()
			txType := "0"
			txSubg := "0"
			txPrnt := "0"
			txData := "8d3729b91a13878508c564fbf410ae4f33fcb4cfdb99677f4b23d4c4adb447650964b4fe9da16299831b9cc17aaabd5b8d81fb05460be92af99d128584101a30" // ?
			txPrev := "c66f4851618cd53104d4a395212958abf88d96962c0c298a0c7a7c1242fac5c2ee616c8c4f140a2e199558ead6d18ae263b2311b590b0d7bf3777be5b3623d9c" // RockSteady was here
			hash := sha512.Sum512([]byte(txTime + txType + txData + txPrev))
			txHash := hex.EncodeToString(hash[:])
			txMile := true
			txLead := false
			txEpoc := txHash
			tx := db.MustBegin()
			tx.MustExec("INSERT INTO transactions (tx_time, tx_type, tx_hash, tx_data, tx_prev, tx_epoc, tx_subg, tx_prnt, tx_mile, tx_lead ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", txTime, txType, txHash, txData, txPrev, txEpoc, txSubg, txPrnt, txMile, txLead)
			tx.Commit()
			return nil
		} else if count > 0 {
			fmt.Printf("%s%s\n"+brightmagenta+"%v"+brightpurple+" transactions known", brightcyan+"\nRoot Tx is ", brightgreen+"present", count)
			return errors.New("Root tx already present. ")
		}
	}
	return nil
}

// createTransaction This will add a transaction to the graph
func createTransaction(txType, data string) {
	var txData string
	var txSubg string
	var txPrev string
	var txPrnt string
	var txLead bool

	// if isCoordinator && txType == "2" {
	if txType == "2" {
		parsePayload := json.Valid([]byte(data))
		if !parsePayload {
			txData = hex.EncodeToString([]byte(data))
		} else if parsePayload {
			txData = data
		}

		db, connectErr := connect()
		defer db.Close()
		handle("Error creating a DB connection: ", connectErr)

		_ = db.QueryRow("SELECT tx_hash FROM transactions ORDER BY tx_time DESC LIMIT 1").Scan(&txPrev)

		txTime := unixTimeStampNano()
		txHash := hashTransaction(txTime, txType, txData, txPrev)
		txEpoc := "0"
		txMile := false
		tx := db.MustBegin()

		if txCount == 0 {
			txLead = true
			txSubg = txHash
			txPrnt = txEpoc
			thisSubgraph = txHash
			txPrnt = thisSubgraph
			thisSubgraphShortName = thisSubgraph[0:4]
			go newSubGraphTimer()
		} else if txCount > 0 {
			txLead = false
			txPrnt = thisSubgraph
			txSubg = thisSubgraph
			thisSubgraphShortName = thisSubgraph[0:4]
		}
		tx.MustExec("INSERT INTO transactions (tx_time, tx_type, tx_hash, tx_data, tx_prev, tx_epoc, tx_subg, tx_prnt, tx_mile, tx_lead ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", txTime, txType, txHash, txData, txPrev, txEpoc, txSubg, txPrnt, txMile, txLead)
		tx.Commit()
		txCount++
	}
}

// newSubGraphTimer timer for collection interval
func newSubGraphTimer() {

	// fmt.Printf(brightcyan+"\nSubgraph created:"+brightgreen+" %s.."+brightcyan+" SubGraph Interval: "+brightgreen+"%vs\n"+nc, thisSubgraph[0:8], poolInterval)
	fmt.Printf(brightcyan+"\nSubgraph created:"+brightgreen+" %s.. "+nc, thisSubgraph[0:8])
	time.Sleep(time.Duration(poolInterval) * time.Second)
	txCount = 0
	// fmt.Printf(brightyellow + "\nInterval concluded" + nc)
}

// hashTransaction takes elements of a transaction and computes a hash using SHA512
func hashTransaction(txTime, txType, txData, txPrev string) string {
	hashedData := []byte(txTime + txType + txData + txPrev)
	slot := make([]byte, 64)
	sha3.ShakeSum256(slot, hashedData)
	// fmt.Printf("%x\n", slot)
	txHash := hex.EncodeToString(slot[:])
	// legacy sha512
	// hash := sha512.Sum512(hashedData)
	// txHash := hex.EncodeToString(hash[:])

	return txHash
}
