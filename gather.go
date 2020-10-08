package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getDataCovid19(seconds time.Duration) {
	dateTime := time.Now()
	fmt.Printf(white+"\n[%s] COVID\t", dateTime.Format("2006-01-02-15:04:05"))
	fmt.Printf(green + "✔️     " + white)
	req, _ := http.NewRequest(
		"GET",
		"https://covidtracking.com/api/v1/us/current.json",
		nil)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	handle("Error: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error: ", err)
	createTransaction("2", string(body))
	delay(seconds)
	go getDataCovid19(seconds)
}

func getDataOgre(seconds time.Duration) {
	dateTime := time.Now()
	fmt.Printf(white+"\n[%s] OGRE\t", dateTime.Format("2006-01-02-15:04:05"))
	fmt.Printf(green + "✔️     " + white)
	req, _ := http.NewRequest(
		"GET",
		"https://tradeogre.com/api/v1/markets",
		nil)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	handle("Error: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error: ", err)
	createTransaction("2", string(body))
	delay(seconds)
	go getDataOgre(seconds)
}
