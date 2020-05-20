package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Hello Karai
func main() {
	parseFlags()
	announce()
	clearPeerID(configPeerIDFile)
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

func sendTransaction(w http.ResponseWriter, r *http.Request) {

}
