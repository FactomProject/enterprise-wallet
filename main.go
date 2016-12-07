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
	"os"
	"os/signal"
	"syscall"

	"github.com/FactomProject/M2GUIWallet/wallet"
	"github.com/FactomProject/factomd/util"
)

var MasterWallet *wallet.WalletDB

func close() {
	fmt.Println("Shutting down")
	if MasterWallet == nil {
		return
	}
	err := MasterWallet.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close()
		os.Exit(1)
	}()
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
	txDB := wallet.MAP     // Holds transactions cache

	fmt.Printf("Starting wallet waspi on localhost:%d\n", walletPort)
	fmt.Printf("Wallet DB using %s, GUI DB using %s, TX DB using %s\n", IntToStringDBType(walletDB), IntToStringDBType(guiDB), IntToStringDBType(txDB))

	// Can adjust starting variables
	// This will also start wallet wsapi
	wal, err := wallet.StartWallet(walletPort, factomdPort, walletDB, guiDB, txDB)
	if err != nil {
		panic("Error in starting wallet: " + err.Error())
	}

	MasterWallet = wal

	// Start Settings
	MasterSettings = new(SettingsStruct)
	data, err := MasterWallet.GUIlDB.Get([]byte("gui-wallet"), []byte("settings"), MasterSettings)
	if err != nil || data == nil {
		err = MasterWallet.GUIlDB.Put([]byte("gui-wallet"), []byte("settings"), MasterSettings)
		if err != nil {
			panic("Error in loading settings: " + err.Error())
		}
	} else {
		MasterSettings = data.(*SettingsStruct)
	}

	// For Testing adds random addresses
	addRandomAddresses()
	//MasterWallet.AddBalancesToAddresses()
	//

	port := 8091
	fmt.Printf("Starting wallet on localhost:%d\n", port)
	ServeWallet(port)
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
