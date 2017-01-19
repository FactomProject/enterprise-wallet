package wallet

import (
	"github.com/FactomProject/factom"
	// "github.com/FactomProject/factom/wallet/wsapi"
)

// StartWallet :
// Must give the port for the factomd instance
func StartWallet(factomdLocation string, walletDBType int, guiDBType int, txDBType int, v1Import bool) (*WalletDB, error) {
	// Set ports
	// factom.SetWalletServer("localhost:" + fmt.Sprintf("%d", walletPort))
	factom.SetFactomdServer(factomdLocation) //"localhost:" + fmt.Sprintf("%d", factomdPort))

	// Can change to MAP, LDB, BOLT
	GUI_DB = guiDBType
	WALLET_DB = walletDBType
	TX_DB = txDBType

	// Load the databases relavent to the wallet
	wal, err := LoadWalletDB(v1Import)
	if err != nil {
		return nil, err
	}

	// This is now uneeded.
	// portStr := "localhost:" + strconv.Itoa(walletPort)
	// fmt.Println("Starting Wallet WSAPI on http://" + portStr + "/")
	// go wsapi.Start(wal.Wallet, fmt.Sprintf(":%d", walletPort), *(factom.RpcConfig))

	return wal, nil
}
