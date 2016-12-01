package TestHelper

import (
	"fmt"

	. "github.com/FactomProject/M2GUIWallet/wallet"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet/wsapi"
)

func Start() (*WalletDB, error) {
	// Should read from config
	factom.SetWalletServer("localhost:8089")
	factom.SetFactomdServer("localhost:8088")

	wal, err := LoadWalletDB()
	if err != nil {
		return nil, err
	}

	// TODO: Adjust start of WSAPI
	go wsapi.Start(wal.Wallet, fmt.Sprintf(":%d", 8089), *(factom.RpcConfig))

	return wal, nil
}
