package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// Graph This is the structure of the Graph
type Graph struct {
	Transactions []*GraphTx `json:"graph_transactions"`
}

// GraphTx This is the structure of the transaction
type GraphTx struct {
	Type int    `json:"tx_type"`
	Hash string `json:"tx_hash"`
	Data string `json:"tx_data"`
	Prev string `json:"tx_prev"`
}

// printTx This is just your basic 'print a transaction'
// command. It takes a file as a parameter.
func printTx(file string) string {
	dat, err := ioutil.ReadFile(file)
	handle("derp, something went wrong", err)
	datString := string(dat) + ",\n"
	return datString
}

// hashTx This will compute the tx hash using sha512
func (graphTx *GraphTx) hashTx() {
	data := strings.Join([]string{graphTx.Data, graphTx.Prev}, "")
	hash := sha512.Sum512([]byte(data))
	// fmt.Printf(brightgreen+"\n%x\n"+white, hash)
	graphTx.Hash = fmt.Sprintf("%x", hash[:])
}

// addTx This will add a transaction to the graph
func (graph *Graph) addTx(txType int, data string) {
	if isCoordinator {
		theTxBeforeThisOne := len(graph.Transactions) - 1
		prevTx := graph.Transactions[theTxBeforeThisOne]
		new := txConstructor(txType, data, prevTx.Hash)
		graph.Transactions = append(graph.Transactions, new)
		if wantsFiles {
			txHeight := len(graph.Transactions) - 1
			heightString := strconv.Itoa(txHeight)
			graphFileName := graphDir + "/" + heightString + ".json"
			data, _ := json.Marshal(new)
			dataToString := string(data)
			createFile(graphFileName)
			writeFile(graphFileName, dataToString)
		}
		if wantsMatrix {
			publishToMatrix(data, matrixURL, matrixRoomID, matrixToken)
		}
	} else {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	}
}

// txConstructor This will construct a tx
func txConstructor(txType int, data, prevHash string) *GraphTx {
	transaction := &GraphTx{txType, "", data, prevHash}
	transaction.hashTx()
	return transaction
}

// rootTx Transaction channels start with a rootTx transaction always
func rootTx() *GraphTx {
	var data = fmt.Sprintf("[%s] Hello Karai", unixTimeStampNano())
	rootName := graphDir + "/0.json"
	thisRoot := txConstructor(0, data, "")
	if !fileExists(rootName) {
		createFile(rootName)
	}
	if wantsFiles {
		thisRootJSON, _ := json.Marshal(thisRoot)
		writeFile(rootName, string(thisRootJSON))
	}

	return thisRoot
}

func prettyPrintGraphJSON(graph *Graph) string {
	graphTransactions := graph.Transactions
	gtJSON, err := json.MarshalIndent(graphTransactions, "", "  ")
	handle("Encountered issues marshalling this JSON: ", err)
	jsonToReturn := fmt.Sprintf("%s\n", string(gtJSON))
	return jsonToReturn
}

// spawnGraph starts a new transaction channel with Root Tx
func spawnGraph() *Graph {
	return &Graph{[]*GraphTx{rootTx()}}
}

// spawnChannel Create a Tx Channel, Root Tx and Milestone, listen for Tx
func spawnChannel() {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		spawnGraph()
	}
}

func writeGraph(graph *Graph) {
	graphJSON := prettyPrintGraphJSON(graph)
	unixFileName := graphDir + "/_graph.json"
	createFile(unixFileName)
	writeFile(unixFileName, graphJSON)
}
func writeTransactions(graph *Graph) {
	start := time.Now()
	graphJSON := graph.Transactions
	chunkSize := 1000000
	customQueue := &workQueue{
		stack: make([]string, 0),
	}
	for i := 0; i < len(graphJSON); i += chunkSize {
		customQueue.push(string(i))
	}
	fmt.Printf(brightcyan+"\nQueue Size: %d"+white, customQueue.size())
	for batch := 0; batch < len(graphJSON); batch += chunkSize {
		fileName := graphDir + "/" + strconv.Itoa((batch/chunkSize)+1) + ".json"
		createFile(fileName)
		data := make([]*GraphTx, chunkSize)
		for item := batch; item < batch+chunkSize; item++ {
			data[item-batch] = graphJSON[item]
		}
		dataJSON, _ := json.Marshal(data)
		writeFile(fileName, string(dataJSON))
		fmt.Printf(brightgreen+"\nSaved: %s"+white, fileName)
		customQueue.pop()
	}
	elapsed := time.Since(start)
	fmt.Printf("\nTook %s seconds.", elapsed)
}
