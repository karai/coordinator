package main

// Hello Karai
func main() {
	osCheck()
	flags()
	checkDirs()
	cleanData()
	keys := initKeys()
	createRoot()
	go restAPI(keys)
	ascii()
	go generateRandomTransactions()
	inputHandler(keys)
}
