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
	p2pListener()
	inputHandler()
}
