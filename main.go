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
	go getDataCovid19(1000)
	go getDataOgre(500)
	go generateRandomTransactions()
	inputHandler(keys)
}
