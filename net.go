package main

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"net"
	"time"
)

func p2pTCPDialer(ip, port, message string, pubKey []byte) {
	var dialer net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	connection, err := dialer.DialContext(ctx, "tcp", ip+":"+port)
	handle("Something went wrong while trying to dial the coordinator: ", err)
	defer connection.Close()
	// TODO add some type of pause here to listen for commands.
	// The way this is currently built means a new connection
	// will be made for every time dialer is called.
	fmt.Fprintf(connection, message)
	status, err := bufio.NewReader(connection).ReadString('\n')
	fmt.Println(status)
}

func p2pListener() {
	listen, err := net.Listen("tcp", ":4201")
	handle("Something went wrong creating a listener: ", err)
	defer listen.Close()
	for {
		listenerConnection, err := listen.Accept()
		handle("Something went wrong accepting a connection: ", err)
		go handleConnection(listenerConnection)
	}
}

func initConnection(pubKey ed25519.PublicKey) {
	joinAddress := "zeus.karai.io"
	joinAddressPort := "4201"
	joinMessage := "JOIN"
	p2pTCPDialer(joinAddress, joinAddressPort, joinMessage, pubKey)
}

func handleConnection(conn net.Conn) {
	_, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Printf("\nPeer disconnected: %s", conn.RemoteAddr())
		conn.Close()
		return
	}
	// bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	// if err != nil {
	// 	fmt.Printf("Peer disconnected: %s", conn.RemoteAddr())
	// 	conn.Close()
	// 	return
	// }
	// message := string(bufferBytes)
	// clientAddr := conn.RemoteAddr().String()
	// response := fmt.Sprintf(message + " from " + clientAddr + "\n")
	// fmt.Println(response)
	// conn.Write([]byte("Sent: " + response))
	// handleConnection(conn)
}
