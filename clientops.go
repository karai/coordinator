package main

import (
	"fmt"
	"net/url"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
)

type clientHeader struct {
	ClientHeaderAppName    string `json:"client_header_app_name"`
	ClientHeaderAppVersion string `json:"client_header_app_version"`
	ClientHeaderPeerID     string `json:"client_header_peer_id"`
	ClientProtocolVersion  string `json:"client_protocol_version"`
}

func connectChannel(ktx string, pubKey ed25519.PublicKey) {

	color.Set(color.FgHiCyan, color.Bold)
	fmt.Printf("\nReceived connection request with ktx %s", ktx)

	color.Set(color.FgHiWhite, color.Bold)
	// connect
	if isCoordinator {
		logrus.Error("this is for clients only.")
	}
	if !isCoordinator {
		// interrupt := make(chan os.Signal, 1)
		// signal.Notify(interrupt, os.Interrupt)
		urlConnection := url.URL{Scheme: "ws", Host: ktx, Path: "/api/v1/channel"}
		color.Set(color.FgHiGreen, color.Bold)
		fmt.Printf("\nConnecting to %s", urlConnection.String())
		conn, _, err := websocket.DefaultDialer.Dial(urlConnection.String(), nil)
		color.Set(color.FgHiRed, color.Bold)
		handle("There was a problem connecting to the channel: ", err)

		msg := "JOIN " + fmt.Sprintf("%x", pubKey)
		fmt.Printf("Msg: %s\n", msg)

		// Initial Connection Sends N1:PK to Coord
		err = conn.WriteMessage(1, []byte(msg))
		// defer conn.Close()
		// done := make(chan struct{})
		// listen for welcome
		// go func() {
		// 	// defer close(done)
		// 	// if we are returning, validate signed N1:S
		// 	// Upon successful connection, submit joinTx
		// 	// if joinTx published, return true on connectChannel() for success
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				color.Set(color.FgHiRed, color.Bold)
				fmt.Println("\nThere was a problem reading this message:", err)
				return
			}
			fmt.Printf("recv: %s", message)
		}
		// }()
	}
}

// The N1 also needs to know how to construct the join message, so I should add that in parallel to clientops.go
// https://discordapp.com/channels/388915017187328002/453726546868305962/719243359440339023
