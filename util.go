package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

// ascii Splash logo. We used to have a package for this
// but it was only for the logo so why not just static-print it?
func ascii() {
	fmt.Printf("\n")
	if isCoordinator {
		fmt.Printf(brightred)
	}
	if !isCoordinator {
		fmt.Printf(brightcyan)
	}
	fmt.Printf("|   _   _  _  .\n")
	fmt.Printf("|( (_| |  (_| |\n")
	fmt.Println(red + semverInfo())
}

// printLicense Print the license for the user
func printLicense() {
	fmt.Printf(brightgreen + "\n" + appName + " v" + semverInfo() + white + " by " + appDev)
	fmt.Println(brightgreen + "\n" + appRepository + "\n" + appURL + "\n")
	fmt.Println(brightwhite + "\nMIT License\nCopyright (c) 2020-2021 RockSteady, TurtleCoin Developers")
	fmt.Println(brightblack + "\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in allcopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.")
	fmt.Println()
}

// menuVersion Print the version string for the user
func menuVersion() {
	fmt.Println(appName + " - v" + semverInfo())
}

// fileExists Does this file exist?
func fileExists(filename string) bool {
	referencedFile, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !referencedFile.IsDir()
}

// directoryMissing Check if a directory has been abducted
func directoryMissing(dirName string) bool {
	src, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			handle("Something went wrong creating the p2p dir: ", err)
		}
		return true
	}
	if src.Mode().IsRegular() {
		fmt.Println(dirName, " already exists as a file.")
		return false
	}
	return false
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

// initPeerLists Check if p2p directory exists, if it does then check for a
// peer file, if it is not there we generate one, then we open it and see if
// it conforms to what we expect, if it does then announce the peer identity.
func initPeerLists() {
	if _, err := os.Stat(p2pConfigDir); os.IsNotExist(err) {
		os.Mkdir(p2pConfigDir, 0700)
	}
	if _, err := os.Stat(p2pWhitelistDir); os.IsNotExist(err) {
		os.Mkdir(p2pWhitelistDir, 0700)
	}
	if _, err := os.Stat(p2pBlacklistDir); os.IsNotExist(err) {
		os.Mkdir(p2pBlacklistDir, 0700)
	}
	if !directoryMissing(p2pConfigDir) {
		if fileExists(configPeerIDFile) {
			peerIdentity := readFile(configPeerIDFile)
			if len(peerIdentity) > 16 {
				fmt.Printf(white + "Machine ID:\t" + brightblack)
				fmt.Printf("%s", peerIdentity)
			} else if len(peerIdentity) < 16 {
				deleteFile(configPeerIDFile)
				generatePeerID()
			}
		} else if !fileExists(configPeerIDFile) {
			fmt.Printf("\nFile " + configPeerIDFile + " does not exist.")
			createFile(configPeerIDFile)
			generatePeerID()
		}
	} else if directoryMissing(p2pConfigDir) {
		fmt.Println("Directory " + p2pConfigDir + " does not exist.")
	}
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
	// fmt.Println(string(text))
	return text
}

// deleteFile Generic file handler
func deleteFile(filename string) {
	os.Remove(filename)
}

// locateGraphDir Find graph storage, create if missing.
func locateGraphDir() {
	if _, err := os.Stat(graphDir); os.IsNotExist(err) {
		err = os.MkdirAll("./graph", 0755)
		handle("Error locating graph directory: ", err)
	}
}

func validJSON(stringToValidate string) bool {
	var jsonString json.RawMessage
	return json.Unmarshal([]byte(stringToValidate), &jsonString) == nil
}
