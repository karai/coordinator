package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

func joinChannel(ktx, pubKey, signedKey, ktxCertFileName string, keyCollection *ED25519Keys) *websocket.Conn {
	// if isCoordinator {
	// 	fmt.Printf(brightred + "\nThis is for nodes running in client mode only.")
	// 	return nil
	// }
	fmt.Printf(brightcyan+"\nConnecting:"+white+" %s", ktx)

	// request a websocket connection
	var conn = requestSocket(ktx, "1")

	// using that connection, attempt to join the channel
	var joinedChannel = joinStatement(conn, pubKey[:64])

	// parse channel messages
	socketMsgParser(ktx, pubKey, signedKey, joinedChannel, keyCollection)

	// return the connection
	return conn
}

func joinStatement(conn *websocket.Conn, pubKey string) *websocket.Conn {

	// new users should send JOIN with the pubkey
	if isFNG {
		joinReq := "JOIN " + pubKey[:64]
		_ = conn.WriteMessage(1, []byte(joinReq))
	}
	// returning users should send RTRN and the signed CA cert
	if !isFNG {
		certString := readFile(selfCertFilePath)
		rtrnReq := "RTRN " + pubKey[:64] + " " + certString
		_ = conn.WriteMessage(1, []byte(rtrnReq))
	}
	return conn
}

func returnMessage(conn *websocket.Conn, pubKey string) *websocket.Conn {
	if !isFNG {
		certString := readFile(selfCertFilePath)
		rtrnReq := "RTRN " + pubKey[:64] + " " + certString
		_ = conn.WriteMessage(1, []byte(rtrnReq))
	}
	return conn
}

func requestSocket(ktx, protocolVersion string) *websocket.Conn {
	urlConnection := url.URL{Scheme: "ws", Host: ktx, Path: "/api/v" + protocolVersion + "/channel"}
	conn, _, err := websocket.DefaultDialer.Dial(urlConnection.String(), nil)
	handle(brightred+"There was a problem connecting to the channel: "+brightcyan, err)
	return conn
}

func socketMsgParser(ktx, pubKey, signedKey string, conn *websocket.Conn, keyCollection *ED25519Keys) {
	_, joinResponse, err := conn.ReadMessage()
	handle("There was a problem reading the socket: ", err)
	if strings.HasPrefix(string(joinResponse), "WCBK") {
		isTrusted = true
		isFNG = false
		fmt.Printf(brightgreen + " ✔️\nConnected!\n" + white)
		fmt.Printf("\nType `"+brightpurple+"send %s <JSON>"+white+"` to send a transaction.\n\n", ktx)
	}
	if strings.Contains(string(joinResponse), string(capkMsg)) {
		convertjoinResponseString := string(joinResponse)
		trimNewLinejoinResponse := strings.TrimRight(convertjoinResponseString, "\n")
		trimCmdPrefix := strings.TrimPrefix(trimNewLinejoinResponse, "CAPK ")
		ncasMsgtring := signKey(keyCollection, trimCmdPrefix[:64])
		composedNcasMsgtring := string(ncasMsg) + " " + ncasMsgtring
		_ = conn.WriteMessage(1, []byte(composedNcasMsgtring))
		_, certResponse, err := conn.ReadMessage()
		isFNG = false
		convertStringcertResponse := string(certResponse) // keys := generateKeys()
		trimNewLinecertResponse := strings.TrimRight(convertStringcertResponse, "\n")
		trimCmdPrefixcertResponse := strings.TrimPrefix(trimNewLinecertResponse, "CERT ")
		handle("There was an error receiving the certificate: ", err)
		ktxCertFileName := p2pConfigDir + "/" + ktx + ".cert"
		createFile(ktxCertFileName)
		writeFile(ktxCertFileName, trimCmdPrefixcertResponse[:192])
		fmt.Printf(brightgreen + "\nCert Name: ")
		fmt.Printf(white+"%s", ktxCertFileName)
		fmt.Printf(brightgreen + "\nCert Body: ")
		fmt.Printf(white+"%s", trimCmdPrefixcertResponse[:192])
	}
}

// Send Takes a data string and a websocket connection
func sendV1Transaction(msg string, conn *websocket.Conn) {
	err := conn.WriteMessage(1, []byte(msg))
	handle("There was a problem sending your transaction ", err)
}
