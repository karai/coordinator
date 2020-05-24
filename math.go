package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/sirupsen/logrus"
	rashedCrypto "github.com/turtlecoin/go-turtlecoin/crypto"
	rashedMnemonic "github.com/turtlecoin/go-turtlecoin/walletbackend/mnemonics"
)

// generatePeerIO uses the machine ID to generate a unique string
func generatePeerID() string {
	logrus.Info("generating peer ID")
	machineID, err := machineid.ProtectedID("1f41d1f36f1f5251f32cfe0f1f924")
	handle("There was a problem generating machine ID: ", err)
	fmt.Println(machineID)
	writeFile(configPeerIDFile, machineID)
	return machineID
}

// checkCreds Locate or create Karai credentials
func checkCreds() {
	if _, err := os.Stat(credentialsFile); err == nil {
		logrus.Debug("Karai Credentials Found!")
	} else {
		logrus.Debug("No Credentials Found! Generating Credentials...")
		rashed25519()
	}
}

// rashed25519 Use TRTL Crypto to generate credentials
// TODO: Replace manually entered JSON
func rashed25519() {
	logrus.Debug("Generating credentials")
	priv, pub, err := rashedCrypto.GenerateKeys()
	seed := rashedMnemonic.PrivateKeyToMnemonic(priv)
	timeUnixNow := strconv.FormatInt(time.Now().Unix(), 10)
	// TODO: Replace manually entered JSON
	logrus.Debug("Writing credentials to file")
	writeFile := []byte("{\n\t\"date_generated\": " + timeUnixNow + ",\n\t\"key_priv\": \"" + hex.EncodeToString(priv[:]) + "\",\n\t\"key_pub\": \"" + hex.EncodeToString(pub[:]) + "\",\n\t\"seed\": \"" + seed + "\"\n}")
	logrus.Debug("Writing main file")
	errWriteFile := ioutil.WriteFile("./"+credentialsFile, writeFile, 0644)
	logrus.Debug(errWriteFile)
	handle("Error writing file: ", err)
	logrus.Debug("Writing backup credential file")
	errWriteBackupFile := ioutil.WriteFile("./."+credentialsFile+"."+timeUnixNow+".backup", writeFile, 0644)
	handle("Error writing file backup: ", err)
	logrus.Debug(errWriteBackupFile)
}

// v4ToHex Convert an ip4 to hex
func v4ToHex(addr string) string {
	ip := net.ParseIP(addr).To4()
	buffer := new(bytes.Buffer)
	for _, s := range ip {
		binary.Write(buffer, binary.BigEndian, uint8(s))
	}
	var dec uint32
	binary.Read(buffer, binary.BigEndian, &dec)
	return fmt.Sprintf("%08x", dec)
}

// portToHex Convert a port to hex
func portToHex(port string) string {
	portNum, _ := strconv.ParseUint(port, 10, 16)
	return fmt.Sprintf("%04x", portNum)
}

// generatePointer create the TRTL <=> Karai pointer
func generatePointer() {
	if !isCoordinator {
		fmt.Println("It looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		logrus.Debug("Creating a new Karai <=> TRTL pointer")
		readerKtxIP := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Karai Coordinator IP: ")
		ktxIP, _ := readerKtxIP.ReadString('\n')
		readerKtxPort := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Karai Coordinator Port: ")
		ktxPort, _ := readerKtxPort.ReadString('\n')
		ip := v4ToHex(strings.TrimRight(ktxIP, "\n"))
		port := portToHex(strings.TrimRight(ktxPort, "\n"))
		fmt.Printf("\nGenerating pointer for %s:%s\n", strings.TrimRight(ktxIP, "\n"), ktxPort)
		fmt.Println("Your pointer is: ")
		fmt.Printf("Hex:\t6b747828%s%s29", ip, port)
		fmt.Println("\nAscii:\tktx(" + strings.TrimRight(ktxIP, "\n") + ":" + strings.TrimRight(ktxPort, "\n") + ")")
	}
}
