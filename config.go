package main

import "github.com/gorilla/websocket"

// Attribution constants
const (
	appName        = "go-karai"
	appDev         = "The TurtleCoin Developers"
	appDescription = appName + " is the Go implementation of the Karai network spec. Karai is a universal blockchain scaling solution for distributed applications."
	appLicense     = "https://choosealicense.com/licenses/mit/"
	appRepository  = "https://github.com/karai/go-karai"
	appURL         = "https://karai.io"
)

// File & folder constants
const (
	configDir         = "./config"
	p2pConfigDir      = "./config/p2p"
	p2pWhitelistDir   = p2pConfigDir + "/whitelist"
	p2pBlacklistDir   = p2pConfigDir + "/blacklist"
	certPathDir       = p2pConfigDir + "/cert"
	certPathSelfDir   = certPathDir + "/self"
	certPathRemote    = certPathDir + "/remote"
	pubKeyFilePath    = certPathSelfDir + "/" + "pub.key"
	privKeyFilePath   = certPathSelfDir + "/" + "priv.key"
	signedKeyFilePath = certPathSelfDir + "/" + "signed.key"
	selfCertFilePath  = certPathSelfDir + "/" + "self.cert"
	currentJSON       = "./config/milestone.json"
	p2pConfigFile     = "peer.id"
)

// Channel values
const (
	channelName string = "⏣ Karai"
	channelDesc string = "This is a general purpose channel."
	channelCont string = "rock@karai.io"
)

// Coordinator values
var (
	nodePubKeySignature []byte
	dbUser              string = "postgres"
	dbName              string = "karai"
	dbSSL               string = "disable"
	joinMsg             []byte = []byte("JOIN")
	ncasMsg             []byte = []byte("NCAS")
	capkMsg             []byte = []byte("CAPK")
	certMsg             []byte = []byte("CERT")
	peerMsg             []byte = []byte("PEER")
	pubkMsg             []byte = []byte("PUBK")
	nsigMsg             []byte = []byte("NSIG")
	sendMsg             []byte = []byte("send")
	rtrnMsg             []byte = []byte("RTRN")
	numTx               int
	consumeData         bool = false
	// isCoordinator       bool   = false
	wantsClean   bool   = false
	graphDir     string = ""
	batchDir     string = ""
	showIP       bool   = false
	chunkSize    int
	karaiAPIPort int
	p2pPeerID    string
	upgrader     = websocket.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}
)

// Client Values
var (
	trimmedPubKey string
	isFNG         = true
	isTrusted     = false
)

// Matrix Values
var (
	wantsMatrix  bool   = false
	wantsFiles   bool   = true
	matrixToken  string = ""
	matrixURL    string = ""
	matrixRoomID string = ""
)

// Subgraph values
var (
	thisSubgraph          string = ""
	thisSubgraphShortName string = ""
	// poolSubLeader         string = ""
	arrangePool  bool = false
	poolInterval int  = 10 // seconds
	txCount      int  = 0
)
