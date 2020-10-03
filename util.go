package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
)

// ascii Splash logo. We used to have a package for this
// but it was only for the logo so why not just static-print it?
func ascii() {
	fmt.Printf("\n\n")
	fmt.Printf(green + "|  |/  / /  /\\  \\ |  |)  ) /  /\\  \\ |  |\n")
	fmt.Printf(brightgreen + "|__|\\__\\/__/¯¯\\__\\|__|\\__\\/__/¯¯\\__\\|__| \n")
	fmt.Printf(brightred + "v" + semverInfo() + white)
	// if isCoordinator {
	fmt.Printf(brightred + " coordinator")
	// }
	// if !isCoordinator {
	// 	fmt.Printf(brightgreen + " client")
	// }

}

// StatsDetail is an object containing strings relevant to the status of a coordinator node.
type StatsDetail struct {
	ChannelName        string `json:"channel_name"`
	ChannelDescription string `json:"channel_description"`
	Version            string `json:"version"`
	ChannelContact     string `json:"channel_contact"`
	PubKeyString       string `json:"pub_key_string"`
	TxObjectsOnDisk    int    `json:"tx_objects_on_disk"`
	GraphUsers         int    `json:"tx_graph_users"`
}

func delay(seconds time.Duration) {
	time.Sleep(seconds * time.Second)
}

// printLicense Print the license for the user
func printLicense() {
	fmt.Printf(brightgreen + "\n" + appName + " v" + semverInfo() + white + " by " + appDev)
	fmt.Printf(brightgreen + "\n" + appRepository + "\n" + appURL + "\n")
	fmt.Printf(brightwhite + "\nMIT License\nCopyright (c) 2020-2021 RockSteady, TurtleCoin Developers")
	fmt.Printf(brightblack + "\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in allcopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.")
	fmt.Printf("\n")
}

// menuVersion Print the version string for the user
func menuVersion() {
	fmt.Printf("%s - v%s\n", appName, semverInfo())
}

// fileExists Does this file exist?
func fileExists(filename string) bool {
	referencedFile, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !referencedFile.IsDir()
}

// fileContainsString This is a utility to see if a string in a file exists.
func fileContainsString(str, filepath string) bool {
	accused, _ := ioutil.ReadFile(filepath)
	isExist, _ := regexp.Match(str, accused)
	return isExist
}

// menuExit Exit the program
func menuExit() {
	os.Exit(0)
}

func timeStamp() string {
	current := time.Now()
	return current.Format("2006-01-02 15:04:05")
}

func unixTimeStampNano() string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	return timestamp
}

// Split helps me split up the args after a command
func Split(r rune) bool {
	return r == ':' || r == '.'
}

func writeTxToDisk(gtxType, gtxHash, gtxData, gtxPrev string) {
	timeNano := unixTimeStampNano()
	txFileName := timeNano + ".json"
	createFile(txFileName)
	txJSONItems := []string{gtxType, gtxHash, gtxData, gtxPrev}
	txJSONObject, _ := json.Marshal(txJSONItems)
	fmt.Printf(white+"\nWriting file...\nFileName: %s\nTransaction Body Object\n%s", txFileName, string(txJSONObject))
	writeFile(txFileName, string(txJSONObject))
}
func createDirIfItDontExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		handle("Could not create directory: ", err)
	}
}

// checkDirs Check if directory exists
func checkDirs() {
	createDirIfItDontExist(graphDir)
	createDirIfItDontExist(batchDir)
	createDirIfItDontExist(configDir)
	createDirIfItDontExist(p2pConfigDir)
	createDirIfItDontExist(p2pWhitelistDir)
	createDirIfItDontExist(p2pBlacklistDir)
	createDirIfItDontExist(certPathDir)
	createDirIfItDontExist(certPathSelfDir)
	createDirIfItDontExist(certPathRemote)
}

// handle Ye Olde Error Handler takes a message and an error code
func handle(msg string, err error) {
	if err != nil {
		fmt.Printf(brightred+"\n%s: %s"+white, msg, err)
	}
}

// createFile Generic file handler
func createFile(filename string) {
	var _, err = os.Stat(filename)
	if os.IsNotExist(err) {
		var file, err = os.Create(filename)
		handle("", err)
		defer file.Close()
	}
}

// writeFile Generic file handler
func writeFile(filename, textToWrite string) {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	handle("", err)
	defer file.Close()
	_, err = file.WriteString(textToWrite)
	err = file.Sync()
	handle("", err)
}

// writeFileBytes Generic file handler
func writeFileBytes(filename string, bytesToWrite []byte) {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	handle("", err)
	defer file.Close()
	_, err = file.Write(bytesToWrite)
	err = file.Sync()
	handle("", err)
}

// readFile Generic file handler
func readFile(filename string) string {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	handle("", err)
	defer file.Close()
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			handle("", err)
			break
		}
	}
	return string(text)
}

func readFileBytes(filename string) []byte {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	handle("", err)
	defer file.Close()
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			handle("", err)
			break
		}
	}
	// fmt.Printf(string(text))
	return text
}

// deleteFile Generic file handler
func deleteFile(filename string) {
	err := os.Remove(filename)
	handle("Problem deleting file: ", err)
}

func validJSON(stringToValidate string) bool {
	// var jsonString json.RawMessage
	// return json.Unmarshal([]byte(stringToValidate), &jsonString) == nil
	return json.Valid([]byte(stringToValidate))
}

func zValidJSON(stringToValidate string) bool {
	// var jsonData map[string]string
	// err := json.Unmarshal([]byte(stringToValidate), &jsonData)
	// if err == nil {
	// 	fmt.Printf("\nJSON is valid")
	// 	return true
	// }
	// fmt.Printf("\nJSON is NOT valid: %s", stringToValidate)
	// return false
	return json.Valid([]byte(stringToValidate))
}

func countFilesOnDisk(directory string) string {
	files, _ := ioutil.ReadDir(directory)
	return strconv.Itoa(len(files))
}

func countWhitelistPeers() int {
	directory := p2pWhitelistDir + "/"
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)
	count := 0
	for range dirFiles {
		count++
	}
	return count
}

func cleanData() {
	if wantsClean {
		// cleanse the whitelist
		directory := p2pWhitelistDir + "/"
		dirRead, _ := os.Open(directory)
		dirFiles, _ := dirRead.Readdir(0)
		for index := range dirFiles {
			fileHere := dirFiles[index]
			nameHere := fileHere.Name()
			fullPath := directory + nameHere
			deleteFile(fullPath)
		}

		// cleanse the blacklist
		blackList, _ := ioutil.ReadDir(p2pBlacklistDir + "/")
		for _, f := range blackList {
			fileToDelete := p2pBlacklistDir + "/" + f.Name()
			fmt.Printf("\nDeleting file: %s", fileToDelete)
			deleteFile(f.Name())
		}
		fmt.Printf(brightyellow+"\nPeers clear: %s"+white, brightgreen+"✔️")

		// cleanse the remote certs
		remoteCert, _ := ioutil.ReadDir(certPathRemote + "/")
		for _, f := range remoteCert {
			fileToDelete := certPathRemote + "/" + f.Name()
			fmt.Printf("\nDeleting file: %s", fileToDelete)
			deleteFile(fileToDelete)
		}
		fmt.Printf(brightyellow+"\nCerts clear: %s"+white, brightgreen+"✔️")

		// cleanse the batches
		batchObjects, _ := ioutil.ReadDir(batchDir + "/")
		for _, f := range batchObjects {
			fileToDelete := batchDir + "/" + f.Name()
			fmt.Printf("\nDeleting file: %s", fileToDelete)
			deleteFile(fileToDelete)
		}
		fmt.Printf(brightyellow+"\nBatches clear: %s"+white, brightgreen+"✔️")

		// cleanse the graph
		graphObjects, _ := ioutil.ReadDir(graphDir + "/")
		for _, f := range graphObjects {
			fileToDelete := graphDir + "/" + f.Name()
			fmt.Printf("\nDeleting file: %s", fileToDelete)
			deleteFile(fileToDelete)
		}
		fmt.Printf(brightyellow+"\nGraph clear: %s"+white, brightgreen+"✔️")
	}
}
