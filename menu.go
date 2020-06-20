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
func inputHandler() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n%v%v%v\n", color.WhiteString("Type '"), color.GreenString("menu"), color.WhiteString("' to view a list of commands"))
		fmt.Print(color.GreenString("-> "))
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
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
			printGraph(graphDir)
		} else if strings.HasPrefix(text, "connect") {
			connectChannel(strings.TrimPrefix(text, "connect "))
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
	color.Set(color.FgGreen)
	fmt.Println("\nCHANNEL_OPTIONS")
	color.Set(color.FgWhite)
	if !isCoordinator {
	} else {
		fmt.Println("create-channel \t\t Create a karai transaction channel")
		fmt.Println("generate-pointer \t Generate a Karai <=> TRTL pointer")
		fmt.Println("benchmark \t\t Conducts timed benchmark")
		fmt.Println("push-graph \t\t Prints graph history")
	}
	color.Set(color.FgGreen)
	fmt.Println("\nWALLET_API_OPTIONS")
	color.Set(color.FgWhite)
	fmt.Println("open-wallet \t\t Open a TRTL wallet")
	fmt.Println("open-wallet-info \t Show wallet and connection info")
	fmt.Println("create-wallet \t\t Create a TRTL wallet")
	color.Set(color.FgHiBlack)
	fmt.Println("wallet-balance \t\t Displays wallet balance")
	color.Set(color.FgGreen)
	fmt.Println("\nKARAI_OPTIONS")
	color.Set(color.FgWhite)
	fmt.Println("connect <ktx> \t Connects to channel where <ktx> is ip.ip.ip.ip:port")
	color.Set(color.FgHiBlack)
	fmt.Println("list-servers \t\t Lists pinning servers")
	color.Set(color.FgGreen)
	fmt.Println("\nGENERAL_OPTIONS")
	color.Set(color.FgWhite)
	fmt.Println("version \t\t Displays version")
	fmt.Println("license \t\t Displays license")
	fmt.Println("exit \t\t\t Quit immediately")
	fmt.Println("")
}
