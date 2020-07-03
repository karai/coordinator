package main

// Hello Karai
func main() {
	parseFlags()
	announce()
	checkPeerFile()
	keys := generateKeys()
	locateGraphDir()
	ascii()
	initAPI(keys)
	go initConnection(keys.publicKey)
	go p2pListener()
	inputHandler(keys)
}
