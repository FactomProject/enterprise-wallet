package TestHelper

import (
	. "github.com/FactomProject/enterprise-wallet/wallet"
	"github.com/FactomProject/factom"
)

func Stop() {
	// Used to do something
}

func Start() (*WalletDB, error) {
	// Should read from config
	factom.SetFactomdServer("localhost:8088")

	wal, err := LoadWalletDB(false, "")
	if err != nil {
		return nil, err
	}

	return wal, nil
}
