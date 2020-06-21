package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// inputHandler This is a basic input loop that listens for
// a few words that correspond to functions in the app. When
// a command isn't understood, it displays the help menu and
// returns to listening to input.
func inputHandler(keyCollection *keys) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n%v%v%v\n", color.WhiteString("Type '"), color.GreenString("menu"), color.WhiteString("' to view a list of commands"))
		fmt.Print(color.GreenString("-> "))
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if strings.Compare("help", text) == 0 {
			menu()
		} else if strings.Compare("?", text) == 0 {
			menu()
		} else if strings.Compare("menu", text) == 0 {
			menu()
		} else if strings.Compare("version", text) == 0 {
			logrus.Debug("Displaying version")
			menuVersion()
		} else if strings.Compare("license", text) == 0 {
			logrus.Debug("Displaying license")
			printLicense()
		} else if strings.Compare("create-wallet", text) == 0 {
			logrus.Debug("Creating Wallet")
			menuCreateWallet()
		} else if strings.Compare("open-wallet", text) == 0 {
			logrus.Debug("Opening Wallet")
			menuOpenWallet()
		} else if strings.Compare("transaction-history", text) == 0 {
			logrus.Debug("Opening Transaction History")
			menuGetContainerTransactions()
		} else if strings.Compare("push-graph", text) == 0 {
			logrus.Debug("Opening Graph History")
			createCID()
		} else if strings.Compare("open-wallet-info", text) == 0 {
			logrus.Debug("Opening Wallet Info")
			menuOpenWalletInfo()
		} else if strings.Compare("benchmark", text) == 0 {
			logrus.Debug("Benchmark")
			benchmark()
		} else if strings.Compare("print-graph", text) == 0 {
			logrus.Debug("Print-graph")
			openGraph(graphDir)
		} else if strings.HasPrefix(text, "connect") {
			connectChannel(strings.TrimPrefix(text, "connect "), keyCollection.pubKey)
		} else if strings.Compare("exit", text) == 0 {
			logrus.Warning("Exiting")
			menuExit()
		} else if strings.Compare("create-channel", text) == 0 {
			logrus.Debug("Creating Karai Transaction Channel")
			spawnChannel()
		} else if strings.Compare("generate-pointer", text) == 0 {
			generatePointer()
		} else if strings.Compare("quit", text) == 0 {
			logrus.Warning("Exiting")
			menuExit()
		} else if strings.Compare("close", text) == 0 {
			logrus.Warning("Exiting")
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
	menuOptions := []string{"CHANNEL_OPTIONS", "WALLET_API_OPTIONS", "KARAI_OPTIONS", "GENERAL_OPTIONS"}
	menuData := map[string][][]string{
		"CHANNEL_OPTIONS": {
			{
				"create-channel \t\t Create a karai transaction channel",
				"generate-pointer \t Generate a Karai <=> TRTL pointer",
				"benchmark \t\t Conducts timed benchmark",
				"push-graph \t\t Prints graph history",
			},
		},
		"WALLET_API_OPTIONS": {
			{
				"open-wallet \t\t Open a TRTL wallet",
				"open-wallet-info \t Show wallet and connection info",
				"create-wallet \t\t Create a TRTL wallet",
			},
			{
				"wallet-balance \t\t Displays wallet balance",
			},
		},
		"KARAI_OPTIONS": {
			{
				"connect <ktx> \t\t Connects to channel where <ktx> is ip.ip.ip.ip:port",
			},
			{
				"list-servers \t\t Lists pinning servers",
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
		color.Set(color.FgGreen)
		fmt.Println("\n" + opt)
		if opt == "CHANNEL_OPTIONS" && !isCoordinator {
			continue
		}
		for colour, options := range menuData[opt] {
			switch colour {
			case 0:
				color.Set(color.FgWhite)
			case 1:
				color.Set(color.FgHiBlack)
			}
			for _, message := range options {
				fmt.Println(message)
			}
		}
	}

	fmt.Println("")
}
