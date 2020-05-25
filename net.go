package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"
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
	logrus.Info("Status: ", status)
}

func p2pListener() {
	// Listen on TCP karaiport on all available unicast and
	// anycast IP addresses of the local system.
	listen, err := net.Listen("tcp", ":4201")
	// listen, err := net.Listen("tcp", ":"+string(karaiP2PPort))
	handle("Something went wrong creating a listener: ", err)
	defer listen.Close()
	for {
		listenerConnection, err := listen.Accept()
		handle("Something went wrong accepting a connection: ", err)
		go func(connection net.Conn) {
			fmt.Println("Got a connection!")
			io.Copy(connection, connection)
			connection.Close()
		}(listenerConnection)
	}
}

func initConnection(ip, port string, pubKey []byte) {
	joinMessage := "JOIN"
	p2pDialer(ip, port, joinMessage, pubKey)
}
