package TestHelper

import (
	"fmt"

	. "github.com/FactomProject/enterprise-wallet/wallet"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet/wsapi"
)

func Stop() {
	wsapi.Stop()
}

func Start(port int) (*WalletDB, error) {
	// Should read from config
	factom.SetWalletServer(fmt.Sprintf("localhost:%d", port))
	factom.SetFactomdServer("localhost:8088")

	wal, err := LoadWalletDB(false)
	if err != nil {
		return nil, err
	}

	// TODO: Adjust start of WSAPI
	go wsapi.Start(wal.Wallet, fmt.Sprintf(":%d", port), *(factom.RpcConfig))
	return wal, nil
}
