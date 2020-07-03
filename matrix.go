package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// publishToMatrix Push events to matrix
func publishToMatrix(matrixMessage, matrixURL, matrixRoomID, matrixToken string) {
	url := matrixURL + "/_matrix/client/r0/rooms/%21" + matrixRoomID + "/send/m.room.message?access_token=" + matrixToken
	// fmt.Printf("\n\nURL GENERATED %s\n", url)
	data := []byte(`{"msgtype":"m.text", "body":"` + matrixMessage + `"}`)
	fmt.Printf("\nDATA GENERATED %s\n\n", data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	handle("Error publishing message to matrix: ", err)
	// req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 4}
	// logrus.Info(req.Header)
	resp, err := client.Do(req)
	handle("Error: ", err)
	defer resp.Body.Close()
	logrus.Info("response Status:", resp.Status)
	logrus.Info("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error: ", err)
	fmt.Printf("%s\n", body)
}
