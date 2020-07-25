package main

// Hello Karai
func main() {
	osCheck()
	flags()
	checkDirs()
	cleanData()
	keys := initKeys()
	graph := initAPI(keys)
	ascii()
	inputHandler(keys, graph)
}
