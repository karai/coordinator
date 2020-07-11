package main

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// Graph This is the structure of the Graph
type Graph struct {
	Transactions []*GraphTx `json:"graph_transactions"`
}

// GraphTx This is the structure of the transaction
type GraphTx struct {
	Type int             `json:"tx_type"`
	Hash string          `json:"tx_hash"`
	Data json.RawMessage `json:"tx_data"`
	Prev string          `json:"tx_prev"`
}

// MilestoneTx This is the structure of the transaction
type MilestoneTx struct {
	Type int    `json:"tx_type"`
	Hash []byte `json:"tx_hash"`
	Data []byte `json:"tx_data"`
	Prev []byte `json:"tx_prev"`
}

// JoinTx This is the structure of the transaction
type JoinTx struct {
	Type int    `json:"tx_type"`
	Hash []byte `json:"tx_hash"`
	Data []byte `json:"tx_data"`
	Prev []byte `json:"tx_prev"`
}

// openGraph a function to open and print the first transaction in graph dir
func openGraph(directory string) {
	handle, err := os.Open(directory + "/" + "Tx_1.json")

	if err != nil {
		fmt.Println("Derp we can't open this JSON: ", err)
	}

	defer handle.Close()
	printGraph(handle)
}

// printGraph a different way to look at transaction history
// this should probably be deleted.
func printGraph(graphHandle io.Reader) {
	var graphTx GraphTx
	if err := json.NewDecoder(graphHandle).Decode(&graphTx); err != nil {
		fmt.Printf("error deserializing JSON: %v", err)
		return
	}

	fmt.Printf("\nhere we go\n%s\n%s\n%d\n%s",
		graphTx.Hash, graphTx.Prev, graphTx.Type, string(graphTx.Data))
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
	data := bytes.Join([][]byte{graphTx.Data, []byte(graphTx.Prev)}, []byte{})
	hash := sha512.Sum512(data)
	graphTx.Hash = string(hash[:])
}

// addTx This will add a transaction to the graph
func (graph *Graph) addTx(txType int, data string) {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		// I wonder sometimes if all these debug statements are costing me tx speed.
		prevTx := graph.Transactions[len(graph.Transactions)-1]
		new := txConstructor(txType, data, []byte(prevTx.Hash))
		graph.Transactions = append(graph.Transactions, new)
		if wantsMatrix {
			publishToMatrix(data, matrixURL, matrixRoomID, matrixToken)
		}

	}
}

// addMilestone This will add a milestone to the graph
func (graph *Graph) addMilestone(data string) {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		prevTransaction := graph.Transactions[len(graph.Transactions)-1]
		// paramFile, _ = os.Open("./config/milestone.json")
		new := txConstructor(1, data, []byte(prevTransaction.Hash))
		graph.Transactions = append(graph.Transactions, new)
		// if wantsMatrix {
		// 	publishToMatrix(data, matrixURL, matrixRoomID, matrixToken)
		// }
	}
}

// txConstructor This will construct a tx
func txConstructor(txType int, data string, prevHash []byte) *GraphTx {
	transaction := &GraphTx{txType, string([]byte{}), []byte(data), string(prevHash)}
	transaction.hashTx()
	return transaction
}

// rootTx Transaction channels start with a rootTx transaction always
func rootTx() *GraphTx {
	fmt.Printf("Coordinator status: %t", isCoordinator)
	var data = "Karai Transaction Channel - Root"
	if wantsMatrix {
		publishToMatrix(data, matrixURL, matrixRoomID, matrixToken)
	}
	return txConstructor(0, data, []byte{})
}

// spawnGraph starts a new transaction channel with Root Tx
func spawnGraph() *Graph {
	return &Graph{[]*GraphTx{rootTx()}}
}

// loadMilestoneJSON Read pending milestone Tx JSON
func loadMilestoneJSON() string {
	// TODO: Check if milestone is ready first, avoid re-use
	dat, _ := ioutil.ReadFile(currentJSON)
	datMilestone := string(dat)
	return datMilestone
	// Kek
}

// validateKTX This function should take a KTX string as a parameter
// and validate that it contains a valid IP and port inside.
// This is used as part of the channel connection process.
func validateKTX(channel string) bool {
	// TODO validate the ktx string with regex
	// if it is valid, return bool true
	return true
}

// sendClientHeader This should batch the client header information
// and send it to the coordinator via the channel parameter. This is
// part of the channel connection process.
func sendClientHeader(name, version, id, channel string) bool {
	// var clientHeaderAppName string = appName
	// var clientHeaderAppVersion string = semverInfo()
	// var clientHeaderPeerID string
	return true
}

// generalHash This is a test function that will probably go away
// soon. It's just a general hash function to hash the milestone
// data returned during the channel connection process.
func (graphTx *GraphTx) generalHash(response string) [64]byte {
	hashedData := bytes.Join([][]byte{graphTx.Data, []byte(graphTx.Prev)}, []byte{})
	hash := sha512.Sum512(hashedData)
	return hash
}

// benchmark Add a number of transactions and time the execution
func benchmark() {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")

	} else {
		benchTxCount := 1000000
		graph := spawnGraph()
		graph.addMilestone(loadMilestoneJSON())
		count := 0
		ascii()
		fmt.Printf("Benchmark: %d transactions\n", benchTxCount)
		fmt.Println("Starting in 5 seconds. Press CTRL C to interrupt.")
		time.Sleep(5 * time.Second)
		start := time.Now()
		for i := 1; i < benchTxCount; i++ {
			count += i
			dataString := "{\"tx_slot\": " + strconv.Itoa(i+1) + "}"
			graph.addTx(2, dataString)
		}
		end := time.Since(start)
		fmt.Printf("\n\nTx Legend: %v %v %v\n", color.YellowString("Root"), color.GreenString("Milestone"), color.BlueString("Normal"))
		for key, transaction := range graph.Transactions {
			var hash string = fmt.Sprintf("%x", transaction.Hash)
			var prevHash string = fmt.Sprintf("%x", transaction.Prev)
			// Root Tx will not have a previous hash
			if prevHash == "" {
				dataString := "{\n\t\"tx_type\": " + strconv.Itoa(transaction.Type) + ",\n\t\"tx_hash\": \"" + hash + "\",\n\t\"tx_data\": \"" + string(transaction.Data) + "\"\n}"
				f, _ := os.Create(graphDir + "/" + "Tx_" + strconv.Itoa(key) + ".json")
				w := bufio.NewWriter(f)
				w.WriteString(dataString)
				w.Flush()
				// fmt.Printf("\nTx(%x) %x\n", key, transaction.Hash)
				fmt.Printf("\nTx(%v) %x\n", color.YellowString(strconv.Itoa(key)), transaction.Hash)
			} else if len(prevHash) > 2 {
				dataString := "{\n\t\"tx_type\": " + strconv.Itoa(transaction.Type) + ",\n\t\"tx_hash\": \"" + hash + "\",\n\t\"tx_prev\": \"" + prevHash + "\",\n\t\"tx_data\": " + string(transaction.Data) + "\n}"
				f, _ := os.Create(graphDir + "/" + "Tx_" + strconv.Itoa(key) + ".json")
				w := bufio.NewWriter(f)
				w.WriteString(dataString)
				w.Flush()
				// Indicate Tx type by color
				if transaction.Type == 0 {
					// Root Tx
					fmt.Printf("Tx(%v) %x\n", color.YellowString(strconv.Itoa(key)), transaction.Hash)
				} else if transaction.Type == 1 {
					// Milestone Tx
					fmt.Printf("Tx(%v) %x\n", color.GreenString(strconv.Itoa(key)), transaction.Hash)
				} else if transaction.Type == 2 {
					// Normal Tx
					fmt.Printf("Tx(%v) %x\n", color.BlueString(strconv.Itoa(key)), transaction.Hash)
				}
			}
		}
		fmt.Println()
		fmt.Printf("%d Transactions in %s", benchTxCount, end)
	}
}

// spawnChannel Create a Tx Channel, Root Tx and Milestone, listen for Tx
func spawnChannel() {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		// Generate Root Tx
		graph := spawnGraph()
		// Add the current milestone.json in config
		graph.addMilestone(loadMilestoneJSON())
		graph.addTx(2, "[{\"tx_slot\": 3}]")
		// go txHandler()
		// Report Txs
		fmt.Printf("\n\nTx Legend: %v %v %v\n", color.YellowString("Root"), color.GreenString("Milestone"), color.BlueString("Normal"))
		for key, transaction := range graph.Transactions {
			var hash string = fmt.Sprintf("%x", transaction.Hash)
			var prevHash string = fmt.Sprintf("%x", transaction.Prev)
			// Root Tx will not have a previous hash
			if prevHash == "" {
				dataString := "{\n\t\"tx_type\": " + strconv.Itoa(transaction.Type) + ",\n\t\"tx_hash\": \"" + hash + "\",\n\t\"tx_data\": \"" + string(transaction.Data) + "\"\n}"
				f, _ := os.Create(graphDir + "/" + "Tx_" + strconv.Itoa(key) + ".json")
				w := bufio.NewWriter(f)
				w.WriteString(dataString)
				w.Flush()
				// fmt.Printf("\nTx(%x) %x\n", key, transaction.Hash)
				fmt.Printf("\nTx(%v) %x\n", color.YellowString(strconv.Itoa(key)), transaction.Hash)
			} else if len(prevHash) > 2 {
				dataString := "{\n\t\"tx_type\": " + strconv.Itoa(transaction.Type) + ",\n\t\"tx_hash\": \"" + hash + "\",\n\t\"tx_prev\": \"" + prevHash + "\",\n\t\"tx_data\": " + string(transaction.Data) + "\n}"
				f, _ := os.Create(graphDir + "/" + "Tx_" + strconv.Itoa(key) + ".json")
				w := bufio.NewWriter(f)
				w.WriteString(dataString)
				w.Flush()
				// Indicate Tx type by color
				if transaction.Type == 0 {
					// Root Tx
					fmt.Printf("Tx(%v) %x\n", color.YellowString(strconv.Itoa(key)), transaction.Hash)
				} else if transaction.Type == 1 {
					// Milestone Tx
					fmt.Printf("Tx(%v) %x\n", color.GreenString(strconv.Itoa(key)), transaction.Hash)
				} else if transaction.Type == 2 {
					// Normal Tx
					fmt.Printf("Tx(%v) %x\n", color.BlueString(strconv.Itoa(key)), transaction.Hash)
				}
			}
		}
		fmt.Println()
	}
}
