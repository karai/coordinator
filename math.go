package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// // generatePeerIO uses the machine ID to generate a unique string
// func generatePeerID() string {
// 	machineID, err := machineid.ProtectedID("1f41d1f36f1f5251f32cfe0f1f924")
// 	handle("There was a problem generating machine ID: ", err)
// 	fmt.Printf(machineID)
// 	writeFile(configPeerIDFile, machineID)
// 	return machineID
// }

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
		fmt.Printf("\nIt looks like you're not a channel coordinator. \n Run Karai with '-coordinator' option to run this command.")
	} else {
		readerKtxIP := bufio.NewReader(os.Stdin)
		fmt.Printf("\nEnter Karai Coordinator IP: ")
		ktxIP, _ := readerKtxIP.ReadString('\n')
		readerKtxPort := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter Karai Coordinator Port: ")
		ktxPort, _ := readerKtxPort.ReadString('\n')
		ip := v4ToHex(strings.TrimRight(ktxIP, "\n"))
		port := portToHex(strings.TrimRight(ktxPort, "\n"))
		fmt.Printf("\nGenerating pointer for %s:%s\n", strings.TrimRight(ktxIP, "\n"), ktxPort)
		fmt.Printf("\nYour pointer is: ")
		fmt.Printf("\nHex:\t6b747828%s%s29", ip, port)
		fmt.Printf("\nAscii:\tktx(" + strings.TrimRight(ktxIP, "\n") + ":" + strings.TrimRight(ktxPort, "\n") + ")")
	}
}
