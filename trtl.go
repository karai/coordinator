package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// menuCreateWallet Create a wallet in the wallet-api container
func menuCreateWallet() {
	url := "http://127.0.0.1:8070/wallet/create"
	data := []byte(`{"daemonHost": "127.0.0.1", "daemonPort": 11898, "filename": "karai-wallet.wallet", "password": "supersecretpassword"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	handle("Error creating wallet: ", err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error creating wallet: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error creating wallet: ", err)
	fmt.Printf("%s\n", body)
}

// menuOpenWallet Open a wallet file
func menuOpenWallet() {
	url := "http://127.0.0.1:8070/wallet/open"
	data := []byte(`{"daemonHost": "127.0.0.1", "daemonPort": 11898, "filename": "karai-wallet.wallet", "password": "supersecretpassword"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	handle("Error opening wallet: ", err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error opening wallet: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error opening wallet: ", err)
	fmt.Printf("%s\n", body)
}

// menuOpenWalletInfo Some basic TRTL API stats
func menuOpenWalletInfo() {
	walletInfoPrimaryAddressBalance()
	getNodeInfo()
	getWalletAPIStatus()
}

// menuGetContainerTransactions Get Wallet-API transactions
func menuGetContainerTransactions() {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/transactions", nil)
	handle("Error getting container transactions: ", err)
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error getting container transactions: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error getting container transactions: ", err)
	fmt.Printf("%s\n", body)
}

// getWalletAPIStatus Get Wallet-API status
func getWalletAPIStatus() {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/status", nil)
	handle("Error getting Wallet-API status: ", err)
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error getting Wallet-API status: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error getting Wallet-API status: ", err)
	fmt.Printf("%s\n", body)
}

// getNodeInfo Get TRTL Node Info
func getNodeInfo() {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/node", nil)
	handle("Error getting node info: ", err)
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error getting node info: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error getting node info: ", err)
	fmt.Printf("%s\n", body)
}

// walletInfoPrimaryAddressBalance Get primary TRTL address balance
func walletInfoPrimaryAddressBalance() {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8070/balances", nil)
	handle("Error getting wallet info primary address: ", err)
	req.Header.Set("X-API-KEY", "pineapples")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	handle("Error getting wallet info primary address: ", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle("Error getting wallet info primary address: ", err)
	fmt.Printf("%s\n", body)
}
