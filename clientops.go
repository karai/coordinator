package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
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

	if isCoordinator {
		color.Set(color.FgHiRed, color.Bold)
		fmt.Println("This is for nodes running in client mode only.")
	}
	if !isCoordinator {
		// request a websocket connection
		var conn = requestSocket(ktx, "1")

		// using that connection, attempt to join the channel
		var joinedChannel = joinChannel(conn, pubKey)

		// validate channel messages
		socketValidation(pubKey, signedKey, joinedChannel)
	}
}

func joinChannel(conn *websocket.Conn, pubKey string) *websocket.Conn {
	// create join message
	joinReq := "JOIN " + pubKey

	// Initial Connection Sends N1:PK to Coord
	_ = conn.WriteMessage(1, []byte(joinReq))
	return conn
}

func requestSocket(ktx, protocolVersion string) *websocket.Conn {
	// Construct a URL for the websocket.
	// For now all that exists is the v1 API, at the /channel
	// endpoint.
	urlConnection := url.URL{Scheme: "ws", Host: ktx, Path: "/api/v" + protocolVersion + "/channel"}

	// Announce the URL we are connecting to
	// color.Set(color.FgHiGreen, color.Bold)
	// fmt.Printf("\nConnecting to %s", urlConnection.String())

	// Make the call to the socket using
	// the URL we composed.
	conn, _, err := websocket.DefaultDialer.Dial(urlConnection.String(), nil)
	color.Set(color.FgHiRed, color.Bold)
	handle("There was a problem connecting to the channel: ", err)
	color.Set(color.FgHiCyan, color.Bold)
	return conn
}

func socketValidation(pubKey, signedKey string, conn *websocket.Conn) {
	_, joinResponse, err := conn.ReadMessage()
	if strings.Contains(string(joinResponse), "Welcome back") {
		fmt.Println("\nConnected")
	} else {
		if len(joinResponse) != 128 {
			color.Set(color.FgHiRed, color.Bold)
			fmt.Println("\nThere seems to be a problem: ", string(joinResponse))
			return
		}
		if err != nil {
			color.Set(color.FgHiRed, color.Bold)
			fmt.Println("\nThere was a problem reading this message:", err)
			return
		}

		// The one issue encountered here is the sig being
		// the wrong length, so lets make sure that is 128
		color.Set(color.FgHiGreen, color.Bold)
		if len(joinResponse) == 128 {
			signature := string(bytes.TrimRight(joinResponse, "\n"))
			color.Set(color.FgHiGreen, color.Bold)
			err = conn.WriteMessage(1, pubkMsg)
			_, readMessageRecvCoordPubKey, _ := conn.ReadMessage()
			coordPubkey := string(bytes.TrimRight(readMessageRecvCoordPubKey, "\n"))
			if verifySignedKey(pubKey, coordPubkey, signature) {
				n1smsg := "NSIG" + signedKey
				color.Set(color.FgHiCyan, color.Bold)
				err = conn.WriteMessage(1, []byte(n1smsg))
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
			fmt.Println("\nSomething is very wrong!")
		}
	}
}
