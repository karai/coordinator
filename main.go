package main

import (
	"flag"
	"net/http"

	"github.com/sirupsen/logrus"
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

// Client Header
var clientHeaderAppName string = appName
var clientHeaderAppVersion string = semverInfo()
var clientHeaderPeerID string

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

func sendTransaction(w http.ResponseWriter, r *http.Request) {

}
