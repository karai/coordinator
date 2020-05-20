package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	externalip "github.com/glendc/go-external-ip"
	"github.com/sirupsen/logrus"
)

func revealIP() string {
	// consensus := externalip.
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	handle("Something went wrong getting the external IP: ", err)
	logrus.Info("External IP: ", ip.String())
	return ip.String()
}

// Splash logo
func ascii() {
	fmt.Printf("\n")
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("|   _   _  _  .\n")
	fmt.Printf("|( (_| |  (_| |\n")
}

// Print the license for the user
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

// Print the version string for the user
func menuVersion() {
	fmt.Println(appName + " - v" + semverInfo())
}

// Exit the program
func menuExit() {
	os.Exit(0)
}

func handle(msg string, err error) {
	if err != nil {
		logrus.Error(msg, err)
	}
}

func parseFlags() {
	flag.IntVar(&karaiPort, "port", 4200, "Port to run Karai Coordinator on.")
	flag.BoolVar(&isCoordinator, "coordinator", false, "Run as coordinator.")
	// flag.StringVar(&karaiPort, "karaiPort", "4200", "Port to run Karai")
	flag.Parse()
}

func announce() {
	if isCoordinator {
		logrus.Info("Coordinator: ", isCoordinator)
		revealIP()

		logrus.Info("Running on port: ", karaiPort)
	} else {
		logrus.Debug("launching as normal user on port: ", karaiPort)
	}
}
