package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
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

func initConnection(pubKey string) {
	joinAddress := "zeus.karai.io"
	joinAddressPort := "4201"
	joinMessage := "JOIN"
	p2pTCPDialer(joinAddress, joinAddressPort, joinMessage, pubKey)
}

func handleConnection(conn net.Conn) {
	_, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Printf("\n[%s] [%s] Peer socket opened!", timeStamp(), conn.RemoteAddr())
		conn.Close()
		return
	}
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Printf("Peer disconnected: %s", conn.RemoteAddr())
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
