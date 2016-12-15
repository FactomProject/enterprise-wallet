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

	"github.com/FactomProject/M2GUIWallet/wallet"
	"github.com/FactomProject/factomd/util"
)

var MasterWallet *wallet.WalletDB

func close() {
	fmt.Println("Shutting down gracefully...")
	if MasterWallet == nil {
		return
	}

	err := MasterWallet.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Complete shut down.")
}

// Initiates and serves the guiwallet. If databases are given, they will be attempted to be loaded
// and will be created if they are not found.
func InitiateWalletAndWeb(guiDBStr string, walDBStr string, txDBStr string, port int, v1Import bool, v1Path string) {
	fmt.Println("--------- Initiating GUIWallet ----------")

	filename := util.ConfigFilename() //file name and path to factomd.conf file
	cfg := util.ReadConfig(filename)

	// Ports
	walletPort := cfg.Wallet.Port
	factomdPort := cfg.App.PortNumber
	controlPanelPort := cfg.App.ControlPanelPort
	if cfg.App.ControlPanelSetting == "disabled" {
		controlPanelPort = -1
	}

	wallet.WalletBoltV1 = v1Path

	var (
		walletDB, guiDB, txDB int
	)

	// DB Types
	switch guiDBStr { // Holds names associated with addresses for gui. Also holds settings
	case "Map":
		guiDB = wallet.MAP
	case "Bolt":
		guiDB = wallet.BOLT
	case "LDB":
		guiDB = wallet.LDB
	}
	switch walDBStr { // WalletDB is DB used by wallet wsapi
	case "Map":
		walletDB = wallet.MAP
	case "Bolt":
		walletDB = wallet.BOLT
	case "LDB":
		walletDB = wallet.LDB
	}
	switch txDBStr { // Holds transactions cache
	case "Map":
		txDB = wallet.MAP
	case "Bolt":
		txDB = wallet.BOLT
	case "LDB":
		txDB = wallet.LDB
	}

	fmt.Printf("Starting wallet waspi on localhost:%d\n", walletPort)
	fmt.Printf("Wallet DB using %s, GUI DB using %s, TX DB using %s\n", IntToStringDBType(walletDB), IntToStringDBType(guiDB), IntToStringDBType(txDB))

	// Can adjust starting variables
	// This will also start wallet wsapi
	wal, err := wallet.StartWallet(walletPort, factomdPort, walletDB, guiDB, txDB, v1Import)
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

	MasterSettings.ControlPanelPort = controlPanelPort
	// For Testing adds random addresses
	if ADD_RANDOM_ADDRESSES {
		addRandomAddresses()
	}
	//

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
