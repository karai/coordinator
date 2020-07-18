package main

import "github.com/gorilla/websocket"

// Attribution constants
const appName = "go-karai"
const appDev = "The TurtleCoin Developers"
const appDescription = appName + " is the Go implementation of the Karai network spec. Karai is a universal blockchain scaling solution for distributed applications."
const appLicense = "https://choosealicense.com/licenses/mit/"
const appRepository = "https://github.com/karai/go-karai"
const appURL = "https://karai.io"

// File & folder constants
const currentJSON = "./config/milestone.json"
const graphDir = "./graph"
const p2pConfigDir = "./config/p2p"
const p2pWhitelistDir = p2pConfigDir + "/whitelist"
const p2pBlacklistDir = p2pConfigDir + "/blacklist"
const p2pConfigFile = "peer.id"
const configPeerIDFile = p2pConfigDir + "/" + "peer.id"
const pubKeyFilePath = p2pConfigDir + "/" + "pub.key"
const privKeyFilePath = p2pConfigDir + "/" + "priv.key"
const signedKeyFilePath = p2pConfigDir + "/" + "signed.key"
const selfCertFilePath = p2pConfigDir + "/" + "self.cert"

// Coordinator values
var nodePubKeySignature []byte
var joinMsg []byte = []byte("JOIN")
var ncasMsg []byte = []byte("NCAS")
var capkMsg []byte = []byte("CAPK")
var certMsg []byte = []byte("CERT")
var peerMsg []byte = []byte("PEER")
var pubkMsg []byte = []byte("PUBK")
var nsigMsg []byte = []byte("NSIG")
var tsxnMsg []byte = []byte("SEND")
var rtrnMsg []byte = []byte("RTRN")
var isCoordinator bool = false
// var wantsHTTPS bool = false
var showIP bool = false
var karaiAPIPort int
var karaiP2PPort int
var p2pPeerID string
var sslDomain = "example.com"
var upgrader = websocket.Upgrader{
	EnableCompression: true,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
}

// Client Values
var trimmedPubKey string
var isFNG = true
var isTrusted = false

// Client Header
var clientHeaderAppName string = appName
var clientHeaderAppVersion string = semverInfo()
var clientHeaderPeerID string

// Matrix Values
var wantsMatrix bool = false
var matrixToken string = ""
var matrixURL string = ""
var matrixRoomID string = ""
