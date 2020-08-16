package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var readyCovid bool
var readyBtc bool

func consume(graph *Graph) {
	if isCoordinator {
		if consumeData {
			go getDataCovid19(5000, graph)
			go getDataBitcoin(120, graph)
			go getDataOgre(120, graph)
		}
	}
}

func getDataCovid19(seconds time.Duration, graph *Graph) {
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
	graph.addTx(2, string(body))
	delay(seconds)
	go getDataCovid19(seconds, graph)
}

func getDataBitcoin(seconds time.Duration, graph *Graph) {
	dateTime := time.Now()
	fmt.Printf(white+"\n[%s] BTC\t", dateTime.Format("2006-01-02-15:04:05"))
	fmt.Printf(green + "✔️     " + white)
	req, _ := http.NewRequest(
		"GET",
		"https://api.coindesk.com/v1/bpi/currentprice.json",
		nil)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	handle("Error: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error: ", err)
	graph.addTx(2, string(body))
	delay(seconds)
	go getDataBitcoin(seconds, graph)
}

func getDataOgre(seconds time.Duration, graph *Graph) {
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
	// i think i should json parse here to remove the slashes
	// but there is no json.Parse in go i dont think
	graph.addTx(2, string(body))
	delay(seconds)
	go getDataOgre(seconds, graph)
}
