package main

// // createCID This will take all of the transactions in the
// // graph directory and call createCIDforTx for each file.
// func createCID() {
// 	start := time.Now()
// 	matches, _ := filepath.Glob(graphDir + "/*.json")
// 	for _, match := range matches {
// 		createCIDforTx(match)
// 	}
// 	end := time.Since(start)
// 	fmt.Printf("\nFinished in: ", end)
// }

// // createCIDforTx This will take a file as a parameter and
// // generate IPFS Content IDs for each file given to it.
// func createCIDforTx(file string) string {
// 	dat, _ := ioutil.ReadFile(file)
// 	fmt.Print(brightblack + string(dat) + "\n")
// 	sh := shell.NewShell("localhost:5001")
// 	cid, err := sh.Add(strings.NewReader(string(dat)))
// 	handle("Something went wrong pushing the tx: ", err)
// 	fmt.Printf(brightgreen+"%v %v\n%v %v", brightyellow+"Tx:", brightgreen+file, brightyellow+"CID: ", brightgreen+cid)
// 	appendGraphCID(cid)
// 	return cid
// }

// appendGraphCID This function will take an IPFS content
// ID as a string and append it to a file containing a list
// of all graph TX's. This is probably not a good idea and
// can be done a different way later when we're more coupled
// with ipfs/libp2p
// func appendGraphCID(cid string) {
// 	if !isCoordinator {
// 		fmt.Printf("\nIt looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
// 	} else {
// 		hashfile, err := os.OpenFile(hashDat,
// 			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 		handle("Something went wrong appending the graph CID: ", err)
// 		defer hashfile.Close()
// 		if fileContainsString(cid, hashDat) {
// 			fmt.Printf("%v", brightred+"\nDuplicate! Skipping...\n")
// 		} else {
// 			hashfile.WriteString(cid + "\n")
// 		}
// 	}
// }
