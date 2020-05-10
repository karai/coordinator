package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	externalip "github.com/glendc/go-external-ip"
	"github.com/gorilla/mux"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	rashedCrypto "github.com/turtlecoin/go-turtlecoin/crypto"
	rashedMnemonic "github.com/turtlecoin/go-turtlecoin/walletbackend/mnemonics"
)

// Attribution constants
const appName = "go-karai"
const appDev = "The TurtleCoin Developers"
const appDescription = appName + " - Karai Transaction Channels"
const appLicense = "https://choosealicense.com/licenses/mit/"
const appRepository = "https://github.com/karai/go-karai"
const appURL = "https://karai.io"

// File & folder constants
const credentialsFile = "private_credentials.karai"
const currentJSON = "./config/milestone.json"
const graphDir = "./graph"
const hashDat = graphDir + "/ipfs-hash-list.dat"
const p2pConfigDir = "./config/p2p"
const configPeerIDFile = p2pConfigDir + "/peer.id"

// Coordinator values
var isCoordinator bool = false
var karaiPort int
var p2pPeerID string

// Version string
func semverInfo() string {
	var majorSemver, minorSemver, patchSemver, wholeString string
	majorSemver = "0"
	minorSemver = "5"
	patchSemver = "2"
	wholeString = majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}

// Graph This is the structure of the Graph
type Graph struct {
	transactions []*GraphTx
}

// GraphTx This is the structure of the transaction
type GraphTx struct {
	TxType   int
	Hash     []byte
	Extra    []byte
	PrevHash []byte
	// TxVer int
}

// // SubGraph This is a struct for Tx wave construction
// type SubGraph struct {
// 	subGraphID       int
// 	timeStamp        int64
// 	milestone        int
// 	transactions     []byte
// 	subgraphChildren int
// 	supgraphOrder    int
// 	subgraphSize     int
// 	subgraphPeers    []byte
// 	waveTip          *GraphTx.Hash
// }

func parseFlags() {
	flag.IntVar(&karaiPort, "port", 4200, "Port to run Karai Coordinator on.")
	flag.BoolVar(&isCoordinator, "coordinator", false, "Run as coordinator.")
	// flag.StringVar(&karaiPort, "karaiPort", "4200", "Port to run Karai")
	flag.Parse()
}

func announce() {
	if isCoordinator {
		logrus.Info("Coordinator: ", isCoordinator)
		revealIP()

		logrus.Info("Running on port: ", karaiPort)
	} else {
		logrus.Debug("launching as normal user on port: ", karaiPort)
	}
}

// Hello Karai
func main() {
	parseFlags()
	announce()
	clearPeerID(configPeerIDFile)
	locateGraphDir()
	checkCreds()
	ascii()
	if !isCoordinator {
		logrus.Debug("isCoordinator == false, skipping webserver deployment")
	} else {
		go restAPI()
	}
	inputHandler()
}

func multiAddr() {

	// construct from a string (err signals parse failure)
	m1, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	handle("", err)
	// construct from bytes (err signals parse failure)
	m2, err := ma.NewMultiaddrBytes(m1.Bytes())
	handle("", err)
	fmt.Println(m2)
	// // true
	// // strings.Equal(m1.String(), "/ip4/127.0.0.1/udp/1234")
	// // strings.Equal(m1.String(), m2.String())
	// bytes.Equal(m1.Bytes(), m2.Bytes())
	// m1.Equal(m2)
	// m2.Equal(m1)

	// // tunneling (en/decap)
	// printer, _ := ma.NewMultiaddr("/ip4/192.168.0.13/tcp/80")
	// proxy, _ := ma.NewMultiaddr("/ip4/10.20.30.40/tcp/443")
	// printerOverProxy := proxy.Encapsulate(printer)
	// proxyAgain := printerOverProxy.Decapsulate(printer)
	// fmt.Println(proxyAgain)
}

func restAPI() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", home).Methods(http.MethodGet)
	api.HandleFunc("/peer", returnPeerID).Methods(http.MethodGet)
	api.HandleFunc("/version", returnVersion).Methods(http.MethodGet)
	api.HandleFunc("/transactions", returnTransactions).Methods(http.MethodGet)
	// api.HandleFunc("", post).Methods(http.MethodPost)
	// api.HandleFunc("", put).Methods(http.MethodPut)
	// api.HandleFunc("", delete).Methods(http.MethodDelete)

	logrus.Error(http.ListenAndServe(":"+strconv.Itoa(karaiPort), r))
}

func revealIP() string {
	// consensus := externalip.
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	handle("Something went wrong getting the external IP: ", err)
	logrus.Info("External IP: ", ip.String())
	return ip.String()
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"bruh": "lol"}`))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Hello " + appName + " v" + semverInfo()))
}

func returnPeerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	peerFile, err := os.OpenFile(configPeerIDFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handle("Something went wrong creating a fresh Peer ID file: ", err)
	defer peerFile.Close()

	fileToRead, err := ioutil.ReadFile(configPeerIDFile)
	// fmt.Println(fileToRead)
	handle("Error: ", err)
	w.Write([]byte("{\"p2p_peer_ID\": \"" + string(fileToRead) + "\"}"))

}

// func printFile(fileToPrint string) string {
// 	file, err := ioutil.ReadFile(configPeerIDFile)
// 	handle("There was a problem reading the peer file", err)
// 	fmt.Print(string(file))
// 	return string(file)
// }

func returnVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"karai_version\": \"" + semverInfo() + "\"}"))
}

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

// Splash logo
func ascii() {
	fmt.Printf("\n")
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("|   _   _  _  .\n")
	fmt.Printf("|( (_| |  (_| |\n")
}

// checkCreds locate or create Karai credentials
func checkCreds() {
	if _, err := os.Stat(credentialsFile); err == nil {
		logrus.Debug("Karai Credentials Found!")
	} else {
		logrus.Debug("No Credentials Found! Generating Credentials...")
		generateEd25519()
	}
}

// generateEd25519 use TRTL Crypto to generate credentials
func generateEd25519() {
	logrus.Debug("Generating credentials")
	priv, pub, err := rashedCrypto.GenerateKeys()
	seed := rashedMnemonic.PrivateKeyToMnemonic(priv)
	timeUnixNow := strconv.FormatInt(time.Now().Unix(), 10)
	// TODO: Replace manually entered JSON
	logrus.Debug("Writing credentials to file")
	writeFile := []byte("{\n\t\"date_generated\": " + timeUnixNow + ",\n\t\"key_priv\": \"" + hex.EncodeToString(priv[:]) + "\",\n\t\"key_pub\": \"" + hex.EncodeToString(pub[:]) + "\",\n\t\"seed\": \"" + seed + "\"\n}")
	logrus.Debug("Writing main file")
	errWriteFile := ioutil.WriteFile("./"+credentialsFile, writeFile, 0644)
	logrus.Debug(errWriteFile)
	handle("Error writing file: ", err)
	logrus.Debug("Writing backup credential file")
	errWriteBackupFile := ioutil.WriteFile("./."+credentialsFile+"."+timeUnixNow+".backup", writeFile, 0644)
	handle("Error writing file backup: ", err)
	logrus.Debug(errWriteBackupFile)
}

// hashTx This will compute the tx hash using sha256
func (graphTx *GraphTx) hashTx() {
	// logrus.Debug("Hashing a Tx ", graphTx.Hash)
	data := bytes.Join([][]byte{graphTx.Extra, graphTx.PrevHash}, []byte{})
	hash := sha256.Sum256(data)
	graphTx.Hash = hash[:]
}

// addTx This will add a transaction to the graph
func (graph *Graph) addTx(txType int, data string) {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		logrus.Debug("Adding a Tx")
		prevTx := graph.transactions[len(graph.transactions)-1]
		new := txConstructor(txType, data, prevTx.Hash)
		graph.transactions = append(graph.transactions, new)
	}
}

func createCID() {
	start := time.Now()
	matches, _ := filepath.Glob(graphDir + "/*.json")
	for _, match := range matches {
		pushTx(match)
	}
	end := time.Since(start)
	fmt.Println("Finished in: ", end)
}

func pushTx(file string) string {
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

func printTx(file string) string {
	dat, _ := ioutil.ReadFile(file)
	datString := string(dat) + ",\n"
	return datString
}

func appendGraphCID(cid string) {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		hashfile, err := os.OpenFile(hashDat,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		handle("Something went wrong appending the graph CID: ", err)
		defer hashfile.Close()
		if isExist(cid, hashDat) {
			fmt.Printf("%v", color.RedString("\nDuplicate! Skipping...\n"))
		} else {
			hashfile.WriteString(cid + "\n")
		}
	}
}

func isExist(str, filepath string) bool {
	accused, _ := ioutil.ReadFile(filepath)
	isExist, _ := regexp.Match(str, accused)
	return isExist
}

// addMilestone This will add a milestone to the graph
func (graph *Graph) addMilestone(data string) {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		prevTransaction := graph.transactions[len(graph.transactions)-1]
		// paramFile, _ = os.Open("./config/milestone.json")
		new := txConstructor(1, data, prevTransaction.Hash)
		graph.transactions = append(graph.transactions, new)
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

// v4ToHex convert an ip4 to hex
func v4ToHex(addr string) string {
	ip := net.ParseIP(addr).To4()
	buffer := new(bytes.Buffer)
	for _, s := range ip {
		binary.Write(buffer, binary.BigEndian, uint8(s))
	}
	var dec uint32
	binary.Read(buffer, binary.BigEndian, &dec)
	return fmt.Sprintf("%08x", dec)
}

// portToHex convert a port to hex
func portToHex(port string) string {
	portNum, _ := strconv.ParseUint(port, 10, 16)
	return fmt.Sprintf("%04x", portNum)
}

// generatePointer create the TRTL <=> Karai pointer
func generatePointer() {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		logrus.Debug("Creating a new Karai <=> TRTL pointer")
		readerKtxIP := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Karai Coordinator IP: ")
		ktxIP, _ := readerKtxIP.ReadString('\n')
		readerKtxPort := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Karai Coordinator Port: ")
		ktxPort, _ := readerKtxPort.ReadString('\n')
		ip := v4ToHex(strings.TrimRight(ktxIP, "\n"))
		port := portToHex(strings.TrimRight(ktxPort, "\n"))
		fmt.Printf("\nGenerating pointer for %s:%s\n", strings.TrimRight(ktxIP, "\n"), ktxPort)
		fmt.Println("Your pointer is: ")
		fmt.Printf("Hex:\t6b747828%s%s29", ip, port)
		fmt.Println("\nAscii:\tktx(" + strings.TrimRight(ktxIP, "\n") + ":" + strings.TrimRight(ktxPort, "\n") + ")")
	}
}

// loadMilestoneJSON Read pending milestone Tx JSON
func loadMilestoneJSON() string {
	// TODO: Check if milestone is ready first, avoid re-use
	dat, _ := ioutil.ReadFile(currentJSON)
	datMilestone := string(dat)
	return datMilestone
	// Kek
}

// // txHandler Wait for Tx, assemble subgraph
// func txHandler() {
// 	var txListenTime time.Duration = 10
// 	var txPoolDepth int = 0
// 	if txPoolDepth > 0 {
// 		// if a tx is received, start the interval, listen for Tx, assemble subgraph
// 		// var int64 SubGraph.timeStamp = time.Now().Unix()
// 		// fmt.Println("Transaction Wave Forming...\nTimestamp: " + string(SubGraph.timeStamp))
// 		time.Sleep(txListenTime * time.Second)
// 		fmt.Println("Listening for " + string(txListenTime) + " seconds")
// 	}
// 	// order the transactions
// 	// assign positions on graph
// }

func connectChannel(channel string) (bool, string) {
	ctx := context.Background()
	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Ping(false),
	)
	handle("Something went wrong creating new peer context: ", err)
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	color.Set(color.FgHiCyan, color.Bold)
	addrs, _ := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("\nlibp2p node address:")
	color.Set(color.FgHiBlack, color.Bold)
	fmt.Println(addrs[0])
	// peerEndPointText, _ := fmt.Println(addrs[0])
	addr, _ := multiaddr.NewMultiaddr(channel)
	peer, _ := peerstore.AddrInfoFromP2pAddr(addr)
	color.Set(color.FgHiCyan, color.Bold)
	fmt.Println("\nConnecting to: ")
	color.Set(color.FgHiBlack, color.Bold)
	fmt.Println(addr)
	ch := pingService.Ping(ctx, peer.ID)
	res := <-ch
	color.Set(color.FgHiRed, color.Bold)
	if err := node.Close(); err != nil {
		handle("\nClosing: ", err)
	}
	if res.RTT >= 0 {
		color.Set(color.FgHiGreen, color.Bold)
		fmt.Println("\nConnected to: ", addr, "in", res.RTT)
		return true, addrs[0].String()
	} else {
		color.Set(color.FgHiRed, color.Bold)
		return false, addrs[0].String()
	}
}

func clearPeerID(file string) {
	err := os.Remove(file)
	logrus.Debug(err)
}

func menuCreatePeer(channel string) {
	// clearPeerID(configPeerIDFile)
	status, p2pPeerID := connectChannel(channel)
	if status == true {
		fmt.Println("Success!")
	} else {
		color.Set(color.FgHiRed, color.Bold)
		fmt.Println("Connection failed: ", status)
	}
	openPeerIDFile, err := os.OpenFile(configPeerIDFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handle("Something went wrong opening the peer ID: ", err)
	defer openPeerIDFile.Close()

	openPeerIDFile.WriteString(p2pPeerID)
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
		for key, transaction := range graph.transactions {
			var hash string = fmt.Sprintf("%x", transaction.Hash)
			var prevHash string = fmt.Sprintf("%x", transaction.PrevHash)
			// Root Tx will not have a previous hash
			if prevHash == "" {
				dataString := "{\n\t\"tx_type\": " + strconv.Itoa(transaction.TxType) + ",\n\t\"tx_hash\": \"" + hash + "\",\n\t\"tx_data\": \"" + string(transaction.Extra) + "\"\n}"
				f, _ := os.Create(graphDir + "/" + "Tx_" + strconv.Itoa(key) + ".json")
				w := bufio.NewWriter(f)
				w.WriteString(dataString)
				w.Flush()
				// fmt.Printf("\nTx(%x) %x\n", key, transaction.Hash)
				fmt.Printf("\nTx(%v) %x\n", color.YellowString(strconv.Itoa(key)), transaction.Hash)
			} else if len(prevHash) > 2 {
				dataString := "{\n\t\"tx_type\": " + strconv.Itoa(transaction.TxType) + ",\n\t\"tx_hash\": \"" + hash + "\",\n\t\"tx_prev\": \"" + prevHash + "\",\n\t\"tx_data\": " + string(transaction.Extra) + "\n}"
				f, _ := os.Create(graphDir + "/" + "Tx_" + strconv.Itoa(key) + ".json")
				w := bufio.NewWriter(f)
				w.WriteString(dataString)
				w.Flush()
				// Indicate Tx type by color
				if transaction.TxType == 0 {
					// Root Tx
					fmt.Printf("Tx(%v) %x\n", color.YellowString(strconv.Itoa(key)), transaction.Hash)
				} else if transaction.TxType == 1 {
					// Milestone Tx
					fmt.Printf("Tx(%v) %x\n", color.GreenString(strconv.Itoa(key)), transaction.Hash)
				} else if transaction.TxType == 2 {
					// Normal Tx
					fmt.Printf("Tx(%v) %x\n", color.BlueString(strconv.Itoa(key)), transaction.Hash)
				}
			}
		}
		fmt.Println()
	}
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
		for key, transaction := range graph.transactions {
			var hash string = fmt.Sprintf("%x", transaction.Hash)
			var prevHash string = fmt.Sprintf("%x", transaction.PrevHash)
			// Root Tx will not have a previous hash
			if prevHash == "" {
				dataString := "{\n\t\"tx_type\": " + strconv.Itoa(transaction.TxType) + ",\n\t\"tx_hash\": \"" + hash + "\",\n\t\"tx_data\": \"" + string(transaction.Extra) + "\"\n}"
				f, _ := os.Create(graphDir + "/" + "Tx_" + strconv.Itoa(key) + ".json")
				w := bufio.NewWriter(f)
				w.WriteString(dataString)
				w.Flush()
				// fmt.Printf("\nTx(%x) %x\n", key, transaction.Hash)
				fmt.Printf("\nTx(%v) %x\n", color.YellowString(strconv.Itoa(key)), transaction.Hash)
			} else if len(prevHash) > 2 {
				dataString := "{\n\t\"tx_type\": " + strconv.Itoa(transaction.TxType) + ",\n\t\"tx_hash\": \"" + hash + "\",\n\t\"tx_prev\": \"" + prevHash + "\",\n\t\"tx_data\": " + string(transaction.Extra) + "\n}"
				f, _ := os.Create(graphDir + "/" + "Tx_" + strconv.Itoa(key) + ".json")
				w := bufio.NewWriter(f)
				w.WriteString(dataString)
				w.Flush()
				// Indicate Tx type by color
				if transaction.TxType == 0 {
					// Root Tx
					fmt.Printf("Tx(%v) %x\n", color.YellowString(strconv.Itoa(key)), transaction.Hash)
				} else if transaction.TxType == 1 {
					// Milestone Tx
					fmt.Printf("Tx(%v) %x\n", color.GreenString(strconv.Itoa(key)), transaction.Hash)
				} else if transaction.TxType == 2 {
					// Normal Tx
					fmt.Printf("Tx(%v) %x\n", color.BlueString(strconv.Itoa(key)), transaction.Hash)
				}
			}
		}
		fmt.Println()
		fmt.Printf("%d Transactions in %s", benchTxCount, end)
	}
}

// locateGraphDir find graph storage, create if missing.
func locateGraphDir() {
	if _, err := os.Stat(graphDir); os.IsNotExist(err) {
		logrus.Debug("Graph directory does not exist.")
		err = os.MkdirAll("./graph", 0755)
		handle("Error locating graph directory: ", err)
	}
}

// inputHandler present menu, accept user input
func inputHandler() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n%v%v%v\n", color.WhiteString("Type '"), color.GreenString("menu"), color.WhiteString("' to view a list of commands"))
		fmt.Print(color.GreenString("-> "))
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("help", text) == 0 {
			menu()
		} else if strings.Compare("?", text) == 0 {
			menu()
		} else if strings.Compare("menu", text) == 0 {
			menu()
		} else if strings.Compare("version", text) == 0 {
			logrus.Debug("Displaying version")
			menuVersion()
		} else if strings.Compare("license", text) == 0 {
			logrus.Debug("Displaying license")
			printLicense()
		} else if strings.Compare("create-wallet", text) == 0 {
			logrus.Debug("Creating Wallet")
			menuCreateWallet()
		} else if strings.Compare("open-wallet", text) == 0 {
			logrus.Debug("Opening Wallet")
			menuOpenWallet()
		} else if strings.Compare("transaction-history", text) == 0 {
			logrus.Debug("Opening Transaction History")
			menuGetContainerTransactions()
		} else if strings.Compare("push-graph", text) == 0 {
			logrus.Debug("Opening Graph History")
			createCID()
		} else if strings.Compare("open-wallet-info", text) == 0 {
			logrus.Debug("Opening Wallet Info")
			menuOpenWalletInfo()
		} else if strings.Compare("benchmark", text) == 0 {
			logrus.Debug("Benchmark")
			benchmark()
			// } else if strings.Compare("create-peer", text) == 0 {
			// 	menuCreatePeer()
		} else if strings.HasPrefix(text, "connect-channel") {
			menuCreatePeer(strings.TrimPrefix(text, "connect-channel "))
		} else if strings.Compare("show-multiaddr", text) == 0 {
			multiAddr()
		} else if strings.Compare("exit", text) == 0 {
			logrus.Warning("Exiting")
			menuExit()
		} else if strings.Compare("create-channel", text) == 0 {
			logrus.Debug("Creating Karai Transaction Channel")
			spawnChannel()
		} else if strings.Compare("generate-pointer", text) == 0 {
			generatePointer()
		} else if strings.Compare("quit", text) == 0 {
			logrus.Warning("Exiting")
			menuExit()
		} else if strings.Compare("close", text) == 0 {
			logrus.Warning("Exiting")
			menuExit()
		} else if strings.Compare("\n", text) == 0 {
			fmt.Println("")
		} else {
			fmt.Println("\nChoose an option from the menu")
			menu()
		}
	}
}

// provide list of commands
func menu() {
	color.Set(color.FgGreen)
	fmt.Println("\nCHANNEL_OPTIONS")
	color.Set(color.FgWhite)
	if !isCoordinator {
	} else {
		fmt.Println("create-channel \t\t Create a karai transaction channel")
		fmt.Println("generate-pointer \t Generate a Karai <=> TRTL pointer")
		fmt.Println("benchmark \t\t Conducts timed benchmark")
	}
	fmt.Println("push-graph \t\t Prints graph history")
	color.Set(color.FgGreen)
	fmt.Println("\nWALLET_API_OPTIONS")
	color.Set(color.FgWhite)
	fmt.Println("open-wallet \t\t Open a TRTL wallet")
	fmt.Println("open-wallet-info \t Show wallet and connection info")
	fmt.Println("create-wallet \t\t Create a TRTL wallet")
	color.Set(color.FgHiBlack)
	fmt.Println("wallet-balance \t\t Displays wallet balance")
	color.Set(color.FgGreen)
	fmt.Println("\nIPFS_OPTIONS")
	color.Set(color.FgWhite)
	fmt.Println("connect-channel <ktx> \t Connects to channel")
	color.Set(color.FgHiBlack)
	if !isCoordinator {
	} else {
		fmt.Println("show-multiaddr \t\t Displays multiaddr address")
	}
	fmt.Println("list-servers \t\t Lists pinning servers")
	color.Set(color.FgGreen)
	fmt.Println("\nGENERAL_OPTIONS")
	color.Set(color.FgWhite)
	fmt.Println("version \t\t Displays version")
	fmt.Println("license \t\t Displays license")
	fmt.Println("exit \t\t\t Quit immediately")
	fmt.Println("")
}

// Some basic TRTL API stats
func menuOpenWalletInfo() {
	walletInfoPrimaryAddressBalance()
	getNodeInfo()
	getWalletAPIStatus()
}

// Get Wallet-API transactions
func menuGetContainerTransactions() {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/transactions", nil)
	handle("Error getting container transactions: ", err)
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error getting container transactions: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error getting container transactions: ", err)
	fmt.Printf("%s\n", body)
}

// Get Wallet-API status
func getWalletAPIStatus() {
	logrus.Info("[Wallet-API Status]")
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/status", nil)
	handle("Error getting Wallet-API status: ", err)
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error getting Wallet-API status: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error getting Wallet-API status: ", err)
	fmt.Printf("%s\n", body)
}

// Get TRTL Node Info
func getNodeInfo() {
	logrus.Info("[Node Info]")
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/node", nil)
	handle("Error getting node info: ", err)
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error getting node info: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error getting node info: ", err)
	fmt.Printf("%s\n", body)
}

// Get primary TRTL address balance
func walletInfoPrimaryAddressBalance() {
	logrus.Info("[Primary Address]")
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/balances", nil)
	handle("Error getting wallet info primary address: ", err)
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error getting wallet info primary address: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error getting wallet info primary address: ", err)
	fmt.Printf("%s\n", body)
}

// Print the license for the user
func printLicense() {
	fmt.Printf(color.GreenString("\n"+appName+" v"+semverInfo()) + color.WhiteString(" by "+appDev))
	color.Set(color.FgGreen)
	fmt.Println("\n" + appRepository + "\n" + appURL + "\n")

	color.Set(color.FgHiWhite)
	fmt.Println("\nMIT License\nCopyright (c) 2020-2021 RockSteady, TurtleCoin Developers")
	color.Set(color.FgHiBlack)
	fmt.Println("\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in allcopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.")
	fmt.Println()
}

// Create a wallet in the wallet-api container
func menuCreateWallet() {
	logrus.Debug("Creating Wallet")
	url := "http://127.0.0.1:8070/wallet/create"
	data := []byte(`{"daemonHost": "127.0.0.1",	"daemonPort": 11898, "filename": "karai-wallet.wallet", "password": "supersecretpassword"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	handle("Error creating wallet: ", err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	logrus.Info(req.Header)
	resp, err := client.Do(req)
	handle("Error creating wallet: ", err)
	defer resp.Body.Close()
	logrus.Info("response Status:", resp.Status)
	logrus.Info("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error creating wallet: ", err)
	fmt.Printf("%s\n", body)
}

// Open a wallet file
func menuOpenWallet() {
	logrus.Debug("Opening Wallet")
	url := "http://127.0.0.1:8070/wallet/open"
	data := []byte(`{"daemonHost": "127.0.0.1",	"daemonPort": 11898, "filename": "karai-wallet.wallet", "password": "supersecretpassword"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	handle("Error opening wallet: ", err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	logrus.Info(req.Header)
	resp, err := client.Do(req)
	handle("Error opening wallet: ", err)
	defer resp.Body.Close()
	logrus.Info("response Status:", resp.Status)
	logrus.Info("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error opening wallet: ", err)
	fmt.Printf("%s\n", body)
}

// Print the version string for the user
func menuVersion() {
	fmt.Println(appName + " - v" + semverInfo())
}

// Exit the program
func menuExit() {
	os.Exit(0)
}

func handle(msg string, err error) {
	if err != nil {
		logrus.Error(msg, err)
	}
}
