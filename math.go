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
