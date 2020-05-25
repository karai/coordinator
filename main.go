package main

// Hello Karai
func main() {
	parseFlags()
	announce()
	generateEd25519()
	checkPeerFile()
	locateGraphDir()
	checkCreds()
	ascii()
	initAPI()
	initConnection()
	go p2pListener()
	inputHandler()
}
