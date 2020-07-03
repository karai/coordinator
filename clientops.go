package main

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type clientHeader struct {
	ClientHeaderAppName    string `json:"client_header_app_name"`
	ClientHeaderAppVersion string `json:"client_header_app_version"`
	ClientHeaderPeerID     string `json:"client_header_peer_id"`
	ClientProtocolVersion  string `json:"client_protocol_version"`
}

func connectChannel(ktx string, pubKey, signedKey string) {
	color.Set(color.FgHiCyan, color.Bold)
	fmt.Printf("\nConnection request with ktx %s", ktx)
	color.Set(color.FgHiWhite, color.Bold)
	// connect
	if isCoordinator {
		logrus.Error("This is for nodes running in client mode only.")
	}
	if !isCoordinator {
		// Construct a URL for the websocket.
		// For now all that exists is the v1 API, at the /channel
		// endpoint.
		urlConnection := url.URL{Scheme: "ws", Host: ktx, Path: "/api/v1/channel"}

		// Announce the URL we are connecting to
		color.Set(color.FgHiGreen, color.Bold)
		fmt.Printf("\nConnecting to %s", urlConnection.String())

		// Make the call to the socket using
		// the URL we composed.
		conn, _, err := websocket.DefaultDialer.Dial(urlConnection.String(), nil)
		color.Set(color.FgHiRed, color.Bold)
		handle("There was a problem connecting to the channel: ", err)
		color.Set(color.FgHiCyan, color.Bold)

		// Craft a message with JOIN as the first word and
		// our nodes pubkey as the second word
		msg := "JOIN " + fmt.Sprintf("%s", pubKey)
		fmt.Printf("\nSending: %s", msg)
		// Initial Connection Sends N1:PK to Coord
		err = conn.WriteMessage(1, []byte(msg))

		// Conditionally validate the response
		for {

			// The response is the coordinator signature of
			// the node public key we just sent it, so here
			// we are telling karai to listen for the response
			// and consider it as the pubkeysig
			_, readMessageRecvPubKeySig, err := conn.ReadMessage()

			if len(readMessageRecvPubKeySig) != 128 {
				color.Set(color.FgHiRed, color.Bold)
				fmt.Println("\nThe Coordinator Public Key Signature we received was not the correct length. \nIt should be 128 characters.")
				return
			}
			// Print some things to help debug
			// fmt.Printf("\n%s\n", readMessageRecvPubKeySig)
			// fmt.Printf("\n%s\n", pubKey)
			if err != nil {
				color.Set(color.FgHiRed, color.Bold)
				fmt.Println("\nThere was a problem reading this message:", err)
				return
			}

			// The one issue encountered here is the sig being
			// the wrong length, so lets make sure that is 128
			color.Set(color.FgHiGreen, color.Bold)
			if len(readMessageRecvPubKeySig) == 128 {
				signature := string(bytes.TrimRight(readMessageRecvPubKeySig, "\n"))
				// Printing the signature for debugging purposes
				fmt.Println("\nCoord Pubkey Signature: " + signature)

				// Write a message to the coordinator requesting
				// the coordinator pubkey. Store it as a var.
				err = conn.WriteMessage(1, pubKeyMsg)
				_, readMessageRecvCoordPubKey, _ := conn.ReadMessage()
				coordPubkey := string(bytes.TrimRight(readMessageRecvCoordPubKey, "\n"))

				// Print the coordinator pubkey signature for debug
				// fmt.Printf("\nCoord Pubkey Signature: %s\n", readMessageRecvCoordPubKey)

				// use ed25519.Verify
				// ed25519.Verify(publicKey ed25519.PublicKey, message []byte, sig []byte)
				fmt.Println("Received Coord Pub Key: " + coordPubkey)
				fmt.Println("Received Coord Signature: " + signature)

				if verifySignedKey(pubKey, coordPubkey, signature) {
					fmt.Println("\nSuccess! The signature from the Channel Coordinator verifies.")

					// Send the signed key for N1s as bytes
					err = conn.WriteMessage(1, []byte(signedKey))

					// Catch the response, this should ve
					_, n1sresponse, _ := conn.ReadMessage()
					hashedSigCertResponse := string(bytes.TrimRight(n1sresponse, "\n"))
					fmt.Println(hashedSigCertResponse)
				}
			} else {
				fmt.Println("\nDrat! It failed..")
			}
		}
	}
}
