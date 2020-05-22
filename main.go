package main

import (
	"github.com/sirupsen/logrus"
)

// Hello Karai
func main() {
	parseFlags()
	announce()
	checkPeerFile()
	locateGraphDir()
	checkCreds()
	ascii()
	if !isCoordinator {
		logrus.Debug("isCoordinator == false, skipping webserver deployment")
	} else {
		go restAPI()
	}
	inputHandler()
}
