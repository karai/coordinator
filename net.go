package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func p2pTCPDialer(ip, port, message string, pubKey string) {
	var dialer net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	connection, _ := dialer.DialContext(ctx, "tcp", ip+":"+port)
	connection.Close()
}

func p2pListener() {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(karaiP2PPort))
	handle("Something went wrong creating a listener: ", err)
	defer listen.Close()
	for {
		listenerConnection, err := listen.Accept()
		handle("Something went wrong accepting a connection: ", err)
		go handleConnection(listenerConnection)
	}
}

func banPeer(peerPubKey string) {
	fmt.Println("Banning peer: " + peerPubKey[:8] + "...")
	whitelist := p2pWhitelistDir + "/" + peerPubKey + ".cert"
	blacklist := p2pBlacklistDir + "/" + peerPubKey + ".cert"
	err := os.Rename(whitelist, blacklist)
	handle("Error banning peer: ", err)
}

func unBanPeer(peerPubKey string) {
	fmt.Println("Unbanning peer: " + peerPubKey[:8] + "...")
	whitelist := p2pWhitelistDir + "/" + peerPubKey + ".cert"
	blacklist := p2pBlacklistDir + "/" + peerPubKey + ".cert"
	err := os.Rename(blacklist, whitelist)
	handle("Error unbanning peer: ", err)
}

func blackList() {
	color.Set(color.FgHiCyan, color.Bold)
	fmt.Println("Displaying banned peers...")
	files, err := ioutil.ReadDir(p2pBlacklistDir)
	handle("There was a problem retrieving the blacklist: ", err)
	color.Set(color.FgHiRed, color.Bold)
	for _, cert := range files {
		certName := cert.Name()
		bannedPeerPubKey := strings.TrimRight(certName, ".cert")
		fmt.Println(bannedPeerPubKey)
	}
	color.Set(color.FgWhite)
}

func clearBlackList() {
	color.Set(color.FgHiCyan, color.Bold)
	fmt.Println("Clearing banned peers...")
	files, err := ioutil.ReadDir(p2pBlacklistDir)
	handle("There was a problem clearing the blacklist: ", err)
	color.Set(color.FgHiRed, color.Bold)
	for _, cert := range files {
		certName := cert.Name()
		bannedPeerPubKey := strings.TrimRight(certName, ".cert")
		// fmt.Println(bannedPeerPubKey)
		unBanPeer(bannedPeerPubKey)
	}
	color.Set(color.FgWhite)
}

func clearPeerList() {
	color.Set(color.FgHiCyan, color.Bold)
	fmt.Println("Purging peer certificates...")

	color.Set(color.FgHiRed, color.Bold)
	directory := p2pWhitelistDir + "/"
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)
	for index := range dirFiles {
		fileHere := dirFiles[index]
		nameHere := fileHere.Name()
		fmt.Println("Purging: ", nameHere)
		fullPath := directory + nameHere
		os.Remove(fullPath)
	}
	color.Set(color.FgHiYellow, color.Bold)
	fmt.Println("Peer list empty!")

}

func whiteList() {
	color.Set(color.FgHiCyan, color.Bold)
	fmt.Println("Displaying peers...")
	files, err := ioutil.ReadDir(p2pWhitelistDir)
	handle("There was a problem retrieving the peerlist: ", err)
	color.Set(color.FgWhite)
	for _, cert := range files {
		certName := cert.Name()
		peerPubKey := strings.TrimRight(certName, ".cert")
		fmt.Println(peerPubKey)
	}
	color.Set(color.FgWhite)
}

func initConnection(pubKey string) {
	joinAddress := "zeus.karai.io"
	joinAddressPort := "4201"
	joinMessage := "JOIN"
	p2pTCPDialer(joinAddress, joinAddressPort, joinMessage, pubKey)
}

func handleConnection(conn net.Conn) {
	_, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		color.Set(color.FgHiBlue, color.Bold)
		fmt.Printf("[%s] [%s] Peer Awakened\n", timeStamp(), conn.RemoteAddr())
		conn.Close()
		color.Set(color.FgWhite)
		return
	}
}
