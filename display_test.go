package main_test

import (
	"net/http/httptest"
	"testing"

	"github.com/FactomProject/enterprise-wallet/TestHelper"
	"github.com/FactomProject/enterprise-wallet/wallet"

	. "github.com/FactomProject/enterprise-wallet"
)

var TestWallet *wallet.WalletDB

func LoadTestWallet(port int) error {
	if TestWallet != nil { // If already instantiated
		return nil
	}

	wallet.GUI_DB = wallet.MAP
	wallet.WALLET_DB = wallet.MAP
	wallet.TX_DB = wallet.MAP

	wal, err := TestHelper.Start(port)
	if err != nil {
		return err
	}

	TestWallet = wal
	TestWallet.Wallet.TXDB().GetAllTXs()
	return nil
}

func TestDisplay(t *testing.T) {
	LoadTestWallet(7089)
	defer TestHelper.Stop()

	MasterWallet = TestWallet
	MasterSettings = new(SettingsStruct)
	InitTemplate()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8091/?request=synced", nil)
	HandleGETRequests(w, r)

	r = httptest.NewRequest("GET", "localhost:8091/?request=addresses-no-bal", nil)
	HandleGETRequests(w, r)

	r = httptest.NewRequest("GET", "localhost:8091/?request=addresses", nil)
	HandleGETRequests(w, r)

	r = httptest.NewRequest("GET", "localhost:8091/?request=balances", nil)
	HandleGETRequests(w, r)

	r = httptest.NewRequest("GET", "localhost:8091/?request=related-transactions", nil)
	HandleGETRequests(w, r)
}
