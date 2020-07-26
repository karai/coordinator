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
	// consume(graph)
	inputHandler(keys, graph)

}
