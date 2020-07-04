package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

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
		// fmt.Printf("\nConnecting to %s", urlConnection.String())

		// Make the call to the socket using
		// the URL we composed.
		conn, _, err := websocket.DefaultDialer.Dial(urlConnection.String(), nil)
		color.Set(color.FgHiRed, color.Bold)
		handle("There was a problem connecting to the channel: ", err)
		color.Set(color.FgHiCyan, color.Bold)

		// Craft a message with JOIN as the first word and
		// our nodes pubkey as the second word
		joinReq := "JOIN " + pubKey
		// fmt.Printf("\nSending: %s", joinReq)
		// Initial Connection Sends N1:PK to Coord
		err = conn.WriteMessage(1, []byte(joinReq))

		// Conditionally validate the response

		// The response is the coordinator signature of
		// the node public key we just sent it, so here
		// we are telling karai to listen for the response
		// and consider it as the pubkeysig
		_, connectionResponse, err := conn.ReadMessage()
		if strings.Contains(string(connectionResponse), "Welcome back") {
			fmt.Println("\nConnected")
			// keep alive?
		} else {
			if len(connectionResponse) != 128 {
				color.Set(color.FgHiRed, color.Bold)
				// fmt.Println("\nThe Coordinator Public Key Signature we received was not the correct length. \nIt should be 128 characters.")
				// fmt.Println("\"" + string(connectionResponse) + "\"" + " is " + string(len(connectionResponse)) + " characters long.")
				fmt.Println("\nThere seems to be a problem: ", string(connectionResponse))
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
			if len(connectionResponse) == 128 {
				signature := string(bytes.TrimRight(connectionResponse, "\n"))
				// Printing the signature for debugging purposes
				// fmt.Printf("\nCoord Pubkey Signature: %s", signature)

				color.Set(color.FgHiCyan, color.Bold)
				// Write a message to the coordinator requesting
				// the coordinator pubkey. Store it as a var.
				// fmt.Printf("\nSending: PUBK request for Coord pubkey...")

				color.Set(color.FgHiGreen, color.Bold)
				err = conn.WriteMessage(1, pubkMsg)
				_, readMessageRecvCoordPubKey, _ := conn.ReadMessage()
				coordPubkey := string(bytes.TrimRight(readMessageRecvCoordPubKey, "\n"))

				// Print the coordinator pubkey signature for debug
				// fmt.Printf("\nCoord Pubkey Signature: %s\n", readMessageRecvCoordPubKey)

				// fmt.Printf("\nReceived Coord PubKey: %s", coordPubkey)
				// fmt.Printf("\nReceived Coord Signature:\t%s", signature)

				if verifySignedKey(pubKey, coordPubkey, signature) {
					// fmt.Println("\nCoordinator signature verified ✔️")

					n1smsg := "NSIG" + signedKey
					// Send the signed key for N1s as bytes

					color.Set(color.FgHiCyan, color.Bold)
					// fmt.Printf("Sending NSIG message to Coordinator...\n%s", n1smsg)
					err = conn.WriteMessage(1, []byte(n1smsg))

					// Catch the response, this should be green
					color.Set(color.FgHiGreen, color.Bold)
					_, n1sresponse, _ := conn.ReadMessage()
					hashedSigCertResponse := bytes.TrimRight(n1sresponse, "\n")
					hashedSigCertResponseNoPrefix := string(bytes.TrimLeft(hashedSigCertResponse, "CERT "))
					// fmt.Println(hashedSigCertResponseNoPrefix)
					if len(hashedSigCertResponseNoPrefix) == 128 {
						fmt.Printf("\n[%s] [%s] Certificate Granted\n", timeStamp(), conn.RemoteAddr())
						color.Set(color.FgHiCyan, color.Bold)
						fmt.Printf("user> ")
						color.Set(color.FgHiBlack, color.Bold)
						fmt.Printf("%s\n", pubKey)
						color.Set(color.FgHiRed, color.Bold)
						fmt.Printf("cert> ")
						color.Set(color.FgHiBlack, color.Bold)
						fmt.Printf("%s\n", hashedSigCertResponseNoPrefix)
						color.Set(color.FgWhite)
					} else {
						fmt.Printf("%v is the wrong size\n%s", len(hashedSigCertResponseNoPrefix), hashedSigCertResponseNoPrefix)
					}
				}
			} else {
				fmt.Println("\nDrat! It failed..")
			}
		}
	}
}
