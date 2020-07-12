package main

// Hello Karai
func main() {
	parseFlags()
	checkPeerFile()
	keys := generateKeys()
	locateGraphDir()
	osCheck()
	ascii()
	initAPI(keys)
	go initConnection(keys.publicKey)
	go p2pListener()
	inputHandler(keys)
}
