package main

// Hello Karai
func main() {
	parseFlags()
	announce()
	keys := generateEd25519()
	checkPeerFile()
	locateGraphDir()
	checkCreds()
	ascii()
	initAPI(keys)
	go initConnection(keys.pubKey)
	go p2pListener()
	inputHandler(keys)
}
