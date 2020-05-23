package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/sirupsen/logrus"
)

// Graph This is the structure of the Graph
type Graph struct {
	Transactions []*GraphTx `json:"graph_transactions"`
}

// GraphTx This is the structure of the transaction
type GraphTx struct {
	Type int    `json:"tx_type"`
	Hash []byte `json:"tx_hash"`
	Data []byte `json:"tx_data"`
	Prev []byte `json:"tx_prev"`
}

// printGraph a different way to look at transaction history
// this should probably be deleted.
func printGraph(directory string) {
	jsonFile, err := os.Open(directory + "/" + "Tx_1.json")
	handle("Derp we can't open this JSON: ", err)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var Graph Graph
	json.Unmarshal(byteValue, &Graph)
	for i := 0; i < 20; i++ {
		fmt.Println("\nhere we go")
		fmt.Println(Graph.Transactions[i].Hash)
		fmt.Println(Graph.Transactions[i].Prev)
		fmt.Println(Graph.Transactions[i].Type)
		fmt.Println(Graph.Transactions[i].Data)
	}
	defer jsonFile.Close()
}

// createCID This will take all of the transactions in the
// graph directory and call createCIDforTx for each file.
func createCID() {
	start := time.Now()
	matches, _ := filepath.Glob(graphDir + "/*.json")
	for _, match := range matches {
		createCIDforTx(match)
	}
	end := time.Since(start)
	fmt.Println("Finished in: ", end)
}

// createCIDforTx This will take a file as a parameter and
// generate IPFS Content IDs for each file given to it.
func createCIDforTx(file string) string {
	dat, _ := ioutil.ReadFile(file)
	color.Set(color.FgBlack, color.Bold)
	fmt.Print(string(dat) + "\n")
	sh := shell.NewShell("localhost:5001")
	cid, err := sh.Add(strings.NewReader(string(dat)))
	handle("Something went wrong pushing the tx: ", err)
	fmt.Printf(color.GreenString("%v %v\n%v %v", color.YellowString("Tx:"), color.GreenString(file), color.YellowString("CID: "), color.GreenString(cid)))
	appendGraphCID(cid)
	return cid
}

// printTv This is just your basic 'print a transaction'
// command. It takes a file as a parameter.
func printTx(file string) string {
	dat, err := ioutil.ReadFile(file)
	handle("derp, something went wrong", err)
	datString := string(dat) + ",\n"
	return datString
}

// appendGraphCID This function will take an IPFS content
// ID as a string and append it to a file containing a list
// of all graph TX's. This is probably not a good idea and
// can be done a different way later when we're more coupled
// with ipfs/libp2p
func appendGraphCID(cid string) {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		hashfile, err := os.OpenFile(hashDat,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		handle("Something went wrong appending the graph CID: ", err)
		defer hashfile.Close()
		if fileContainsString(cid, hashDat) {
			fmt.Printf("%v", color.RedString("\nDuplicate! Skipping...\n"))
		} else {
			hashfile.WriteString(cid + "\n")
		}
	}
}

// hashTx This will compute the tx hash using sha256
func (graphTx *GraphTx) hashTx() {
	// logrus.Debug("Hashing a Tx ", graphTx.Hash)
	data := bytes.Join([][]byte{graphTx.Data, graphTx.Prev}, []byte{})
	hash := sha256.Sum256(data)
	graphTx.Hash = hash[:]
}

// addTx This will add a transaction to the graph
func (graph *Graph) addTx(txType int, data string) {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		logrus.Debug("Adding a Tx")
		prevTx := graph.Transactions[len(graph.Transactions)-1]
		new := txConstructor(txType, data, prevTx.Hash)
		graph.Transactions = append(graph.Transactions, new)
	}
}

// addMilestone This will add a milestone to the graph
func (graph *Graph) addMilestone(data string) {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		prevTransaction := graph.Transactions[len(graph.Transactions)-1]
		// paramFile, _ = os.Open("./config/milestone.json")
		new := txConstructor(1, data, prevTransaction.Hash)
		graph.Transactions = append(graph.Transactions, new)
	}
}

// txConstructor This will construct a tx
func txConstructor(txType int, data string, prevHash []byte) *GraphTx {
	transaction := &GraphTx{txType, []byte{}, []byte(data), prevHash}
	transaction.hashTx()
	return transaction
}

// rootTx Transaction channels start with a rootTx transaction always
func rootTx() *GraphTx {
	fmt.Printf("Coordinator status: %t", isCoordinator)
	return txConstructor(0, "Karai Transaction Channel - Root", []byte{})
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
func (graphTx *GraphTx) generalHash(response string) [32]byte {
	hashedData := bytes.Join([][]byte{graphTx.Data, graphTx.Prev}, []byte{})
	hash := sha256.Sum256(hashedData)
	return hash
}

// func connectToChannel(channel string) {
//  if validateKTX(channel) {
//      if validateCoordVersion(channel) {
//          //send client header to coord
//          if sendClientHeader(clientHeaderAppName, clientHeaderAppVersion, clientHeaderPeerID, channel) {
//              //coord should respond with most recent milestone
//              //hash the milestone
//              if generalHash(res.Body) == milestone.Hash {
//                  sendClientMilestoneHash(channel)
//              }
//              //send the hash to coord
//              //coord approves
//              //send join tx
//              //listen for events
//          } else if sendClientHeader(clientHeaderAppName, clientHeaderAppVersion, clientHeaderPeerID, channel) {
//              logrus.Error("Problem constructing or sending client header.")
//          }
//      } else if !validateCoordVersion(channel) {
//          logrus.Error("Coordinator Version Not Accepted")
//      }
//  } else if !validateKTX(channel) {
//      logrus.Error("KTX Invalid")
//  }
// }

// func coordVersionHandler(w http.ResponseWriter, r *http.Request) {
//  return
// }

// func validateCoordVersion(channel string) bool {
//  logrus.Info("fetching coordinator version for ", channel)
//  req, err := http.NewRequest("GET", channel, nil)
//  handle("Error getting coord info: ", err)
//  client := &http.Client{Timeout: time.Second * 10}
//  resp, err := client.Do(req)
//  handle("Error getting coord info: ", err)
//  defer resp.Body.Close()
//  body, err := ioutil.ReadAll(resp.Body)
//  logrus.Debug(body)
//  handle("Error getting coord info: ", err)
//  return true
// }

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

// locateGraphDir Find graph storage, create if missing.
func locateGraphDir() {
	if _, err := os.Stat(graphDir); os.IsNotExist(err) {
		logrus.Debug("Graph directory does not exist.")
		err = os.MkdirAll("./graph", 0755)
		handle("Error locating graph directory: ", err)
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
		graph.addTx(2, "{\"tx_slot\": 3}")
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
