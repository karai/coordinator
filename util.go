package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/fatih/color"
	externalip "github.com/glendc/go-external-ip"
	"github.com/sirupsen/logrus"
)

// revealIP This uses some funky consensus methods to
// dial a few servers and get the external IP of the coordinator
func revealIP() string {
	// consensus := externalip.
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	handle("Something went wrong getting the external IP: ", err)
	logrus.Info("External IP: ", ip.String())
	return ip.String()
}

// ascii Splash logo. We used to have a package for this
// but it was only for the logo so why not just static-print it?
func ascii() {
	fmt.Printf("\n")
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("|   _   _  _  .\n")
	fmt.Printf("|( (_| |  (_| |\n")
}

// printLicense Print the license for the user
func printLicense() {
	fmt.Printf(color.GreenString("\n"+appName+" v"+semverInfo()) + color.WhiteString(" by "+appDev))
	color.Set(color.FgGreen)
	fmt.Println("\n" + appRepository + "\n" + appURL + "\n")

	color.Set(color.FgHiWhite)
	fmt.Println("\nMIT License\nCopyright (c) 2020-2021 RockSteady, TurtleCoin Developers")
	color.Set(color.FgHiBlack)
	fmt.Println("\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in allcopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.")
	fmt.Println()
}

// menuVersion Print the version string for the user
func menuVersion() {
	fmt.Println(appName + " - v" + semverInfo())
}

func fileExists(filename string) bool {
	referencedFile, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !referencedFile.IsDir()
}

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

func checkPeerFile() {
	// logrus.Info("Checking peer file: " + p2pConfigDir + "/" + p2pConfigFile)
	if !directoryMissing(p2pConfigDir) {
		// logrus.Info(p2pConfigDir + " exists")
		if fileExists(configPeerIDFile) {
			// logrus.Info(configPeerIDFile + " exists")
			peerIdentity := readFile(configPeerIDFile)
			if len(peerIdentity) > 16 {
				logrus.Debug("Peer identity looks ok")
				logrus.Info("Peer identity: " + peerIdentity)
				// TODO ADD MORE VALIDATION
			} else if len(peerIdentity) < 16 {
				deleteFile(configPeerIDFile)
				generatePeerID()
			}
		} else if !fileExists(configPeerIDFile) {
			logrus.Warning("File " + configPeerIDFile + " does not exist.")
			createFile(configPeerIDFile)
			generatePeerID()
		}
	} else if directoryMissing(p2pConfigDir) {
		fmt.Println("Directory " + p2pConfigDir + " does not exist.")
	}

	// if fileExists(configPeerIDFile) {
	// 	fmt.Println("exists.")
	// 	peerFile, err := ioutil.ReadFile(configPeerIDFile)
	// 	handle("There was a problem reading the peer file: ", err)
	// 	if len(peerFile) > 7 {
	// 		logrus.Info("Peer Identity: ", peerFile)
	// 	} else if !fileExists(configPeerIDFile) {
	// 		fmt.Println("Doesnt Exist.")
	// 		// if we have no peer id in the file or a malformed peer ID,
	// 		// we should generate a new one.
	// 		_, err := os.Stat(p2pConfigDir)
	// 		if os.IsNotExist(err) {
	// 			errDir := os.MkdirAll(p2pConfigDir, 0755)
	// 			if errDir != nil {
	// 				handle("Something went wrong creating the P2P directory: ", err)
	// 			}
	// 		}
	// 		logrus.Warning("Peer Identity is not present. Generating new...")
	// 		generatePeerID()
	// 	}
	// }
}

// handle Ye Olde Error Handler takes a message and an error code
func handle(msg string, err error) {
	if err != nil {
		logrus.Error(msg, err)
	}
}

// parseFlags This evaluates the flags used when the program was run
// and assigns the values of those flags according to sane defaults.
func parseFlags() {
	flag.IntVar(&karaiPort, "port", 4200, "Port to run Karai Coordinator on.")
	flag.BoolVar(&isCoordinator, "coordinator", false, "Run as coordinator.")
	// flag.StringVar(&karaiPort, "karaiPort", "4200", "Port to run Karai")
	flag.Parse()
}

// announce Tell us when the program is running
func announce() {
	if isCoordinator {
		logrus.Info("Coordinator: ", isCoordinator)
		revealIP()

		logrus.Info("Running on port: ", karaiPort)
	} else {
		logrus.Debug("launching as normal user on port: ", karaiPort)
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
	logrus.Debug("Created file: ", filename)
}

// writeFile Generic file handler
func writeFile(filename, textToWrite string) {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	handle("", err)
	defer file.Close()
	_, err = file.WriteString(textToWrite)
	err = file.Sync()
	handle("", err)
	logrus.Debug("Text written to file: ", textToWrite)
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
	logrus.Debug("Read from file: ", text)
	// fmt.Println(string(text))
	return string(text)
}

// deleteFile Generic file handler
func deleteFile(filename string) {
	var err = os.Remove(filename)
	handle("", err)

	logrus.Debug("Deleted file: ", filename)
}
