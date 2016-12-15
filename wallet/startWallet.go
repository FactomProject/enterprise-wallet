package wallet

import (
	"fmt"
	"strconv"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet/wsapi"
)

// Must give the port for the factomd instance and wallet wsapi
// will start wallet wsapi on selected port
func StartWallet(walletPort int, factomdPort int, guiDBType int, walletDBType int, txDBType int, v1Import bool) (*WalletDB, error) {
	// Set ports
	factom.SetWalletServer("localhost:" + fmt.Sprintf("%d", walletPort))
	factom.SetFactomdServer("localhost:" + fmt.Sprintf("%d", factomdPort))

	// Can change to MAP, LDB, BOLT
	GUI_DB = guiDBType
	WALLET_DB = walletDBType
	TX_DB = txDBType

	wal, err := LoadWalletDB(v1Import)
	if err != nil {
		return nil, err
	}

	// TODO: Adjust start of WSAPI -- RpcConfig
	portStr := "localhost:" + strconv.Itoa(walletPort)
	fmt.Println("Starting Wallet WSAPI on http://localhost" + portStr + "/")
	go wsapi.Start(wal.Wallet, fmt.Sprintf(":%d", walletPort), *(factom.RpcConfig))

	return wal, nil
}
