package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type clientHeader struct {
	ClientHeaderAppName    string `json:"client_header_app_name"`
	ClientHeaderAppVersion string `json:"client_header_app_version"`
	ClientHeaderPeerID     string `json:"client_header_peer_id"`
	ClientProtocolVersion  string `json:"client_protocol_version"`
}

func connectChannel(ktx string) {
	// connect
	if !isCoordinator {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		u := url.URL{Scheme: "ws", Host: ktx, Path: "/api/v1/channel"}
		log.Printf("connecting to %s", u.String())
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		handle("There was a problem connecting to the channel: ", err)
		// Initial Connection Sends N1:PK to Coord
		err = conn.WriteMessage(1, []byte("JOIN "+string(pubKey)))
		defer conn.Close()
		done := make(chan struct{})
		// listen for welcome
		go func() {
			defer close(done)
			// if we are returning, validate signed N1:S
			// Upon successful connection, submit joinTx
			// if joinTx published, return true on connectChannel() for success
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("There was a problem reading this message:", err)
					return
				}
				fmt.Printf("recv: %s", message)
			}
		}()
	}
}

// The N1 also needs to know how to construct the join message, so I should add that in parallel to clientops.go
// https://discordapp.com/channels/388915017187328002/453726546868305962/719243359440339023
