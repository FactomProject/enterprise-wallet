package main

/*
 * Begins all the services required by the GUI wallet
 * 		- WSAPI for wallet
 *		- Webserver
 * Requires for all functionality
 *		- Factomd Instance
 */

import (
	"fmt"

	"github.com/FactomProject/M2WalletGUI/wallet"
	"github.com/FactomProject/factomd/util"
)

var MasterWallet *wallet.WalletDB

func main() {
	InitiateWalletAndWeb()
}

func InitiateWalletAndWeb() {
	fmt.Println("--------- Initiating GUIWallet ----------")

	filename := util.ConfigFilename() //file name and path to factomd.conf file
	cfg := util.ReadConfig(filename)

	// Ports
	walletPort := cfg.Wallet.Port
	factomdPort := cfg.App.PortNumber

	// DB Types
	walletDB := wallet.MAP // WalletDB is DB used by wallet wsapi
	guiDB := wallet.MAP    // Holds names associated with addresses for gui

	fmt.Printf("Starting wallet waspi on localhost:%d\n", walletPort)
	fmt.Printf("Wallet DB using %s, GUI DB using %s\n", IntToStringDBType(walletDB), IntToStringDBType(guiDB))

	// Can adjust starting variables
	// This will also start wallet wsapi
	wal, err := wallet.StartWallet(walletPort, factomdPort, walletDB, guiDB)
	if err != nil {
		panic("Error in starting wallet: " + err.Error())
	}

	MasterWallet = wal

	// For Testing
	addRandomAddresses()
	MasterWallet.AddBalancesToAddresses()
	//

	port := 8091
	fmt.Printf("Starting wallet on localhost:%d\n", port)
	ServeWallet(port)

	// Exiting, need to make an exit condition, just a placeholder and reminder for now
	wal.Close()
}

func addRandomAddresses() {
	for i := 0; i < 5; i++ {
		MasterWallet.GenerateEntryCreditAddress("AddedForTesting")
	}

	for i := 0; i < 5; i++ {
		MasterWallet.GenerateFactoidAddress("AddedForTesting")
	}

	MasterWallet.AddAddress("Sand", "Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK")
}

func IntToStringDBType(t int) string {
	switch t {
	case wallet.MAP:
		return "Map"
	case wallet.LDB:
		return "LDB"
	case wallet.BOLT:
		return "Bolt"
	}
	return "[No DB Type Found]"
}
