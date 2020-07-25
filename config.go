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
	currentJSON       = "./config/milestone.json"
	graphDir          = "./graph"
	p2pConfigDir      = "./config/p2p"
	p2pWhitelistDir   = p2pConfigDir + "/whitelist"
	p2pBlacklistDir   = p2pConfigDir + "/blacklist"
	p2pConfigFile     = "peer.id"
	certPath          = p2pConfigDir + "/cert"
	certPathSelf      = certPath + "/self"
	certPathRemote    = certPath + "/remote"
	pubKeyFilePath    = certPathSelf + "/" + "pub.key"
	privKeyFilePath   = certPathSelf + "/" + "priv.key"
	signedKeyFilePath = certPathSelf + "/" + "signed.key"
	selfCertFilePath  = certPathSelf + "/" + "self.cert"
)

// Coordinator values
var (
	nodePubKeySignature []byte
	joinMsg             []byte = []byte("JOIN")
	ncasMsg             []byte = []byte("NCAS")
	capkMsg             []byte = []byte("CAPK")
	certMsg             []byte = []byte("CERT")
	peerMsg             []byte = []byte("PEER")
	pubkMsg             []byte = []byte("PUBK")
	nsigMsg             []byte = []byte("NSIG")
	sendMsg             []byte = []byte("send")
	rtrnMsg             []byte = []byte("RTRN")
	isCoordinator       bool   = false
	wantsClean          bool   = false
)

var (
	showIP       bool = false
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

// Client Header
var (
	clientHeaderAppName    string = appName
	clientHeaderAppVersion string = semverInfo()
	clientHeaderPeerID     string
)

// Matrix Values
var (
	wantsMatrix  bool   = false
	wantsFiles   bool   = true
	matrixToken  string = ""
	matrixURL    string = ""
	matrixRoomID string = ""
)
