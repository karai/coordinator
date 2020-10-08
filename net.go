package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func banPeer(peerPubKey string) {
	fmt.Printf("\nBanning peer: %s" + peerPubKey[:8] + "...")
	whitelist := p2pWhitelistDir + "/" + peerPubKey + ".cert"
	blacklist := p2pBlacklistDir + "/" + peerPubKey + ".cert"
	err := os.Rename(whitelist, blacklist)
	handle("Error banning peer: ", err)
}

func unBanPeer(peerPubKey string) {
	fmt.Printf("\nUnbanning peer: %s" + peerPubKey[:8] + "...")
	whitelist := p2pWhitelistDir + "/" + peerPubKey + ".cert"
	blacklist := p2pBlacklistDir + "/" + peerPubKey + ".cert"
	err := os.Rename(blacklist, whitelist)
	handle("Error unbanning peer: ", err)
}

func blackList() {
	fmt.Printf(brightcyan + "\nDisplaying banned peers...")
	files, err := ioutil.ReadDir(p2pBlacklistDir)
	handle(brightred+"There was a problem retrieving the blacklist: ", err)
	for _, cert := range files {
		certName := cert.Name()
		bannedPeerPubKey := strings.TrimRight(certName, ".cert")
		fmt.Printf(brightred + "\n" + bannedPeerPubKey + white)
	}
}

func clearBlackList() {
	fmt.Printf(brightcyan + "Clearing banned peers..." + white)
	files, err := ioutil.ReadDir(p2pBlacklistDir)
	handle(brightred+"There was a problem clearing the blacklist: ", err)
	for _, cert := range files {
		certName := cert.Name()
		bannedPeerPubKey := strings.TrimRight(certName, ".cert")
		unBanPeer(bannedPeerPubKey)
	}
}

func clearPeerList() {
	fmt.Printf(brightcyan + "Purging peer certificates...")
	directory := p2pWhitelistDir + "/"
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)
	for index := range dirFiles {
		fileHere := dirFiles[index]
		nameHere := fileHere.Name()
		fmt.Printf(brightred+"Purging: %s"+white, nameHere)
		fullPath := directory + nameHere
		os.Remove(fullPath)
	}
	fmt.Printf(brightyellow + "\n" + "Peer list empty!" + white)

}

func whiteList() {
	fmt.Printf(brightcyan + "Displaying peers...\n")
	files, err := ioutil.ReadDir(p2pWhitelistDir)
	handle(brightred+"There was a problem retrieving the peerlist: "+white, err)
	for _, cert := range files {
		certName := cert.Name()
		peerPubKey := strings.TrimRight(certName, ".cert")
		fmt.Printf("\n" + peerPubKey)
	}
}

func handleConnection(conn net.Conn) {
	_, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Printf(brightcyan+"[%s] [%s] Peer Awakened\n"+white, timeStamp(), conn.RemoteAddr())
		conn.Close()
		return
	}
}
