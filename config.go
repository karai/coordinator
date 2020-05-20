package main

// Attribution constants
const appName = "go-karai"
const appDev = "The TurtleCoin Developers"
const appDescription = appName + " is the Go implementation of the Karai network spec. Karai is a universal blockchain scaling solution for distributed applications."
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
