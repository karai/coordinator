package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func p2pDialer(ip, port, message string, pubKey []byte) {
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

func initConnection() {
	joinAddress := "zeus.karai.io"
	joinAddressPort := "4201"
	joinMessage := "JOIN"
	p2pDialer(joinAddress, joinAddressPort, joinMessage, pubKey)
}

// func handleRequest(conn net.Conn) {
// 	buf := make([]byte, 1024)
// 	readBuffer, err := conn.Read(buf)
// 	if err != nil {
// 		fmt.Println("Error reading:", err.Error())
// 	}
// 	fmt.Println(readBuffer)
// 	ackstring := "Message received." + string(signedPubKey)
// 	conn.Write([]byte(ackstring))
// 	conn.Close()
// }

func handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("client left..")
		conn.Close()
		return
	}
	message := string(bufferBytes)
	clientAddr := conn.RemoteAddr().String()
	response := fmt.Sprintf(message + " from " + clientAddr + "\n")
	fmt.Println(response)
	conn.Write([]byte("Sent: " + response))
	handleConnection(conn)
}
