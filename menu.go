package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Split helps me split up the args after a command
func Split(r rune) bool {
	return r == ':' || r == '.'
}

type v1Tx struct {
	channel string
	msg     string
}

// inputHandler This is a basic input loop that listens for
// a few words that correspond to functions in the app. When
// a command isn't understood, it displays the help menu and
// returns to listening to input.
func inputHandler(keyCollection *ED25519Keys) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n%v%v%v\n", white+"Type '", brightgreen+"menu", white+"' to view a list of commands")
		fmt.Print(white + "-> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if strings.Compare("help", text) == 0 {
			menu()
		} else if strings.Compare("?", text) == 0 {
			menu()
		} else if strings.Compare("menu", text) == 0 {
			menu()
		} else if strings.Compare("version", text) == 0 {
			menuVersion()
		} else if strings.Compare("license", text) == 0 {
			printLicense()
		} else if strings.Compare("create-wallet", text) == 0 {
			menuCreateWallet()
		} else if strings.Compare("open-wallet", text) == 0 {
			menuOpenWallet()
		} else if strings.Compare("transaction-history", text) == 0 {
			menuGetContainerTransactions()
		} else if strings.Compare("push-graph", text) == 0 {
			createCID()
		} else if strings.Compare("open-wallet-info", text) == 0 {
			menuOpenWalletInfo()
		} else if strings.Compare("benchmark", text) == 0 {
			benchmark()
		} else if strings.Compare("print-graph", text) == 0 {
			openGraph(graphDir)
		} else if strings.HasPrefix(text, "connect") {
			ktxAddressString := strings.TrimPrefix(text, "connect ")
			if strings.Contains(ktxAddressString, ":") {
				var justTheDomainPartNotThePort = strings.Split(ktxAddressString, ":")
				var ktxCertFileName = justTheDomainPartNotThePort[0] + ".cert"
				if !fileExists(ktxCertFileName) {
					joinChannel(ktxAddressString, keyCollection.publicKey, keyCollection.signedKey, "")
				}
				if fileExists(ktxCertFileName) {
					isFNG = false
					joinChannel(ktxAddressString, keyCollection.publicKey, keyCollection.signedKey, ktxCertFileName)
				}

			}
			if !strings.Contains(ktxAddressString, ":") {
				fmt.Printf("\nDid you forget to include the port?\n")
			}
		} else if strings.HasPrefix(text, "send") {
			// strip away the `send` command prefix
			input := strings.TrimPrefix(text, "send ")

			// split up the words after the prefix
			args := strings.FieldsFunc(input, Split)

			// form an array conforming to our v1 tx type
			tx := v1Tx{args[0], args[1]}

			// extend some data to vars
			channel := fmt.Sprintf(tx.channel)
			msg := fmt.Sprintf(tx.msg)

			// send the v1 transaction if the json is valid
			if validJSON(msg) {
				sendV1Transaction(channel, msg)
			} else {
				fmt.Printf("\nYour transaction %s was not interpreted as valid JSON. ", msg)
			}
		} else if strings.HasPrefix(text, "ban ") {
			bannedPeer := strings.TrimPrefix(text, "ban ")
			banPeer(bannedPeer)
		} else if strings.HasPrefix(text, "unban ") {
			unBannedPeer := strings.TrimPrefix(text, "unban ")
			unBanPeer(unBannedPeer)
		} else if strings.HasPrefix(text, "blacklist") {
			blackList()
		} else if strings.HasPrefix(text, "clear blacklist") {
			clearBlackList()
		} else if strings.HasPrefix(text, "clear peerlist") {
			clearPeerList()
		} else if strings.HasPrefix(text, "peerlist") {
			whiteList()
		} else if strings.Compare("exit", text) == 0 {
			menuExit()
		} else if strings.Compare("create-channel", text) == 0 {
			spawnChannel()
		} else if strings.Compare("generate-pointer", text) == 0 {
			generatePointer()
		} else if strings.Compare("quit", text) == 0 {
			menuExit()
		} else if strings.Compare("close", text) == 0 {
			menuExit()
		} else if strings.Compare("\n", text) == 0 {
			fmt.Println("")
		} else {
			fmt.Println("\nChoose an option from the menu")
			menu()
		}
	}
}

// menu This is the body of text printed when the user
// types 'help', 'menu' or any undefined input.
func menu() {
	menuOptions := []string{"LAUNCH_PARAMETERS", "CHANNEL_OPTIONS", "WALLET_API_OPTIONS", "KARAI_OPTIONS", "GENERAL_OPTIONS"}
	menuData := map[string][][]string{
		"LAUNCH_PARAMETERS": {
			{
				"-coordinator \t\t Run Karai as Coordinator",
				"-https \t\t\t Use HTTPS for Coordinator API",
				"-matrix \t\t Send event messages to Matrix homeserver",
				"-matrixtoken \t\t Matrix homeserver token string",
				"-matrixurl \t\t Matrix homeserver URL string",
				"-matrixroomid \t\t Room ID string for matrix publishd events",
				"-apiport \t\t Coordinator API port integer",
			},
			{},
		},
		"CHANNEL_OPTIONS": {
			{
				"create-channel \t\t Create a karai transaction channel",
				"generate-pointer \t Generate a Karai <=> TRTL pointer",
				"benchmark \t\t Conducts timed benchmark",
				"push-graph \t\t Prints graph history",
			},
			{},
		},
		"WALLET_API_OPTIONS": {
			{},
			{
				"open-wallet \t\t Open a TRTL wallet",
				"open-wallet-info \t Show wallet and connection info",
				"create-wallet \t\t Create a TRTL wallet",
				"wallet-balance \t\t Displays wallet balance",
			},
		},
		"KARAI_OPTIONS": {
			{
				"connect <ktx> \t\t Connects to channel where <ktx> is ip.ip.ip.ip:port",
				"peerlist \t\t Lists known peers.",
				"blacklist \t\t Lists banned peers.",
				"clear blacklist \t Unbans all blacklist peer certificates.",
				"clear peerlist \t\t Purges all whitelist peer certificates.",
				"ban <pubkey> \t\t Ban user certificate by pubkey.",
				"unban <pubkey> \t\t Unban user certificate by pubkey.",
			},
		},
		"GENERAL_OPTIONS": {
			{
				"version \t\t Displays version",
				"license \t\t Displays license",
				"exit \t\t\t Quit immediately",
			},
		},
	}

	for _, opt := range menuOptions {
		fmt.Println(brightgreen + "\n" + opt)
		for menuOptionColor, options := range menuData[opt] {
			switch menuOptionColor {
			case 0:
				fmt.Printf(brightwhite)
			case 1:
				fmt.Printf(brightblack)
			}
			for _, message := range options {
				fmt.Println(message)
			}
		}
	}

	fmt.Println("")
}
