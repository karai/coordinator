package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
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
	if fileExists(rootName) {
		fmt.Printf(brightred + "A root transaction exists. Create-channel refused.")
		fmt.Printf(brightcyan+"Filename: %s"+white, rootName)
	}
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
	graphJSON := graph.Transactions
	for _, object := range graphJSON {
		timeNano := unixTimeStampNano()
		fileName := graphDir + "/" + timeNano + ".json"
		createFile(fileName)
		objJSON, err := json.Marshal(object)
		handle("There were issues marshalling that JSON: ", err)
		writeFile(fileName, string(objJSON))
		fmt.Printf(brightgreen+"\nCreated: %s", fileName)
	}

}

// Commenting this because it needs to be refactored
// At the least, there need to be more fields to access
// attributes of a milestone tx.
// // MilestoneTx This is the structure of the transaction
// type MilestoneTx struct {
// 	Type int    `json:"tx_type"`
// 	Hash string `json:"tx_hash"`
// 	Data string `json:"tx_data"`
// 	Prev string `json:"tx_prev"`
// }

// // addMilestone This will add a milestone to the graph
// func (graph *Graph) addMilestone(data string) {
// 	if !isCoordinator {
// 		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
// 	} else {
// 		prevTransaction := graph.Transactions[len(graph.Transactions)-1]
// 		new := txConstructor(1, data, prevTransaction.Hash)
// 		graph.Transactions = append(graph.Transactions, new)
// 		if wantsMatrix {
// 			publishToMatrix(data, matrixURL, matrixRoomID, matrixToken)
// 		}
// 	}
// }

// // printGraph a different way to look at transaction history
// // this should probably be deleted.
// func printGraph(graphHandle io.Reader) {
// 	var graphTx GraphTx
// 	if err := json.NewDecoder(graphHandle).Decode(&graphTx); err != nil {
// 		fmt.Printf("error deserializing JSON: %v", err)
// 		return
// 	}

// 	fmt.Printf("\nhere we go\n%s\n%s\n%d\n%s",
// 		graphTx.Hash, graphTx.Prev, graphTx.Type, string(graphTx.Data))
// }

// // loadMilestoneJSON Read pending milestone Tx JSON
// func loadMilestoneJSON() string {
// 	// TODO: Check if milestone is ready first, avoid re-use
// 	dat, _ := ioutil.ReadFile(currentJSON)
// 	datMilestone := string(dat)
// 	return datMilestone
// 	// Kek
// }
