package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

// func p2pTCPDialer(ip, port, message string, pubKey string) {
// 	var dialer net.Dialer
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
// 	defer cancel()
// 	connection, _ := dialer.DialContext(ctx, "tcp", ip+":"+port)
// 	connection.Close()
// }

// func p2pListener() {
// 	listen, err := net.Listen("tcp", ":"+strconv.Itoa(karaiP2PPort))
// 	handle("Something went wrong creating a listener: ", err)
// 	defer listen.Close()
// 	for {
// 		listenerConnection, err := listen.Accept()
// 		handle("Something went wrong accepting a connection: ", err)
// 		go handleConnection(listenerConnection)
// 	}
// }

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
	fmt.Println(brightcyan + "Displaying banned peers...")
	files, err := ioutil.ReadDir(p2pBlacklistDir)
	handle(brightred+"There was a problem retrieving the blacklist: ", err)
	for _, cert := range files {
		certName := cert.Name()
		bannedPeerPubKey := strings.TrimRight(certName, ".cert")
		fmt.Println(brightred + bannedPeerPubKey + white)
	}
}

func clearBlackList() {
	fmt.Println(brightcyan + "Clearing banned peers..." + white)
	files, err := ioutil.ReadDir(p2pBlacklistDir)
	handle(brightred+"There was a problem clearing the blacklist: ", err)
	for _, cert := range files {
		certName := cert.Name()
		bannedPeerPubKey := strings.TrimRight(certName, ".cert")
		unBanPeer(bannedPeerPubKey)
	}
}

func clearPeerList() {
	fmt.Println(brightcyan + "Purging peer certificates...")
	directory := p2pWhitelistDir + "/"
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)
	for index := range dirFiles {
		fileHere := dirFiles[index]
		nameHere := fileHere.Name()
		fmt.Println(brightred+"Purging: "+white, nameHere)
		fullPath := directory + nameHere
		os.Remove(fullPath)
	}
	fmt.Println(brightyellow + "Peer list empty!" + white)

}

func whiteList() {
	fmt.Println(brightcyan + "Displaying peers...")
	files, err := ioutil.ReadDir(p2pWhitelistDir)
	handle(brightred+"There was a problem retrieving the peerlist: "+white, err)
	for _, cert := range files {
		certName := cert.Name()
		peerPubKey := strings.TrimRight(certName, ".cert")
		fmt.Println(peerPubKey)
	}
}

// func initConnection(pubKey string) {
// 	joinAddress := "zeus.karai.io"
// 	joinAddressPort := "4201"
// 	joinMessage := "HELLO " + pubKey
// 	p2pTCPDialer(joinAddress, joinAddressPort, joinMessage, pubKey)
// }

func handleConnection(conn net.Conn) {
	_, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Printf(brightcyan+"[%s] [%s] Peer Awakened\n"+white, timeStamp(), conn.RemoteAddr())
		conn.Close()
		return
	}
}
