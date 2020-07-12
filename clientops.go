package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

type clientHeader struct {
	ClientHeaderAppName    string `json:"client_header_app_name"`
	ClientHeaderAppVersion string `json:"client_header_app_version"`
	ClientHeaderPeerID     string `json:"client_header_peer_id"`
	ClientProtocolVersion  string `json:"client_protocol_version"`
}

func joinChannel(ktx string, pubKey, signedKey, ktxCertFileName string) *websocket.Conn {
	fmt.Printf(brightcyan+"\nConnection request with ktx %s", ktx)

	if isCoordinator {
		fmt.Println(brightred + "This is for nodes running in client mode only.")
		return nil
	}

	// request a websocket connection
	var conn = requestSocket(ktx, "1")

	// using that connection, attempt to join the channel
	var joinedChannel = sendJoinMsg(conn, pubKey, ktxCertFileName)

	// parse channel messages
	socketMsgParser(ktx, pubKey, signedKey, joinedChannel)

	// return the connection
	return conn
}

func sendJoinMsg(conn *websocket.Conn, pubKey, ktxCertFileName string) *websocket.Conn {
	// create join message
	if isFNG {
		joinReq := "JOIN " + pubKey
		// Initial Connection Sends N1:PK to Coord
		_ = conn.WriteMessage(1, []byte(joinReq))
	}
	if !isFNG {
		certString := readFile(ktxCertFileName)
		rtrnReq := "RTRN " + pubKey + " " + certString
		// Connection Sends CA:N1:S to Coord
		_ = conn.WriteMessage(1, []byte(rtrnReq))
	}
	return conn
}

func requestSocket(ktx, protocolVersion string) *websocket.Conn {
	// Construct a URL for the websocket.
	// For now all that exists is the v1 API, at the /channel
	// endpoint.
	urlConnection := url.URL{Scheme: "ws", Host: ktx, Path: "/api/v" + protocolVersion + "/channel"}

	// Announce the URL we are connecting to
	// fmt.Printf("\nConnecting to %s", urlConnection.String())

	// Make the call to the socket using
	// the URL we composed.
	conn, _, err := websocket.DefaultDialer.Dial(urlConnection.String(), nil)
	handle(brightred+"There was a problem connecting to the channel: "+brightcyan, err)
	return conn
}

func socketMsgParser(ktx, pubKey, signedKey string, conn *websocket.Conn) {
	_, joinResponse, err := conn.ReadMessage()
	if strings.Contains(string(joinResponse), "Welcome back") {
		fmt.Println("\nConnected")
	} else {
		if len(joinResponse) != 128 {
			fmt.Println(brightred+"\nThere seems to be a problem: ", string(joinResponse))
			return
		}
		if err != nil {
			fmt.Println(brightred+"\nThere was a problem reading this message:", err)
			return
		}

		// The one issue encountered here is the sig being
		// the wrong length, so lets make sure that is 128
		if len(joinResponse) == 128 {
			signature := string(bytes.TrimRight(joinResponse, "\n"))
			err = conn.WriteMessage(1, pubkMsg)
			_, readMessageRecvCoordPubKey, _ := conn.ReadMessage()
			coordPubkey := string(bytes.TrimRight(readMessageRecvCoordPubKey, "\n"))
			if verifySignedKey(pubKey, coordPubkey, signature) {
				n1smsg := "NSIG" + signedKey
				err = conn.WriteMessage(1, []byte(n1smsg))
				_, n1sresponse, _ := conn.ReadMessage()
				hashedSigCertResponse := bytes.TrimRight(n1sresponse, "\n")
				hashedSigCertResponseNoPrefix := string(bytes.TrimLeft(hashedSigCertResponse, "CERT "))
				if len(hashedSigCertResponseNoPrefix) == 128 {
					var justTheDomainPartNotThePort = strings.Split(ktx, ":")
					var ktxCertFileName = justTheDomainPartNotThePort[0] + ".cert"
					if !fileExists(ktxCertFileName) {
						createFile(ktxCertFileName)
					}
					writeFile(ktxCertFileName, hashedSigCertResponseNoPrefix)
					fmt.Printf(brightgreen+"\n[%s] [%s] Certificate Granted\n", timeStamp(), conn.RemoteAddr())
					fmt.Printf("file> ")
					fmt.Printf(brightblack+"./%s", ktxCertFileName)
					fmt.Printf(brightcyan + "\nuser> ")
					fmt.Printf(brightblack+"%s\n", pubKey)
					fmt.Printf(brightred + "cert> ")
					fmt.Printf(brightblack+"%s\n"+white, hashedSigCertResponseNoPrefix)
				} else {
					fmt.Printf("%v is the wrong size\n%s", len(hashedSigCertResponseNoPrefix), hashedSigCertResponseNoPrefix)
				}
			}
		} else {
			fmt.Println("\nSomething is very wrong!")
		}
	}
}

func sendV1Transaction(channel, msg string) {
	conn := requestSocket(channel, "1")
	_ = conn.WriteMessage(1, []byte(msg))

}
