package main

// Hello Karai
func main() {
	parseFlags()
	initPeerLists()
	keys := initKeys()
	locateGraphDir()
	osCheck()
	ascii()
	initAPI(keys)
	// go initConnection(keys.publicKey)
	// go p2pListener()
	inputHandler(keys)
}
