package wallet_test

import (
	"fmt"
	"testing"

	"github.com/FactomProject/enterprise-wallet/TestHelper"
	. "github.com/FactomProject/enterprise-wallet/wallet"
)

var _ = fmt.Sprint("")

var TestWallet *WalletDB

func TestSendFactoids(t *testing.T) {
	//fmt.Println(3)
	var err error
	err = LoadTestWallet(8077)
	if err != nil {
		t.Fatal(err.Error())
	}

	//FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q
	anp, list := TestWallet.GetGUIAddress("FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q")
	if list == -1 {
		anp, err = TestWallet.AddAddress("Sand", "Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK")
		if err != nil {
			t.Fatal(err)
		}
	}

	var recs []string
	var amts []uint64
	for i := 0; i < 20; i++ {
		recAnp, err := TestWallet.GenerateFactoidAddress("SendToThis")
		if err != nil {
			t.Fatal(err, i)
		}

		recs = append(recs, recAnp.Address)
		//amts = append(amts, 1000e8)
		amts = append(amts, 1e8) //1e8)
	}

	/*
		balance, err := TestWallet.GetAddressBalance(recAnp.Address)
		if err != nil {
			t.Fatal(err)
		}
	*/

	//20000e8
	trans, _, err := TestWallet.ConstructTransaction(recs, amts)
	if err != nil {
		t.Fatal(err)
	}

	var total uint64
	var amtsStrs []string
	for _, a := range amts {
		total += a
		amtsStrs = append(amtsStrs, fmt.Sprintf("%d", a/1e8))
	}

	nameComp, err := TestWallet.CheckTransactionAndGetName(recs, amtsStrs, "")
	if err != nil {
		t.Fatal(err)
	} else if trans != nameComp {
		t.Fatal("Names do not match")
	}

	_, err = TestWallet.SendTransaction(trans)
	if err != nil {
		t.Fatal(err)
	}

	// Test string versions
	name, ret, err := TestWallet.ConstructSendFactoidsStrings(recs, amtsStrs)
	if err != nil {
		t.Fatal(err)
	}

	if ret.Total != total {
		t.Fatal("Total wrong")
	}

	_, err = TestWallet.SendTransaction(name)
	if err != nil {
		t.Fatal(err)
	}

	_ = anp
}

func TestConvertToEC(t *testing.T) {
	LoadTestWallet(8076)
	var err error

	//FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q
	anp, list := TestWallet.GetGUIAddress("FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q")
	if list == -1 {
		anp, err = TestWallet.AddAddress("Sand", "Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK")
		if err != nil {
			t.Fatal(err)
		}
	}

	var recs []string
	var amts []uint64
	for i := 0; i < 20; i++ {
		recAnp, err := TestWallet.GenerateEntryCreditAddress("SendToThis")
		if err != nil {
			t.Fatal(err, i)
		}

		recs = append(recs, recAnp.Address)
		//amts = append(amts, 1000e8)
		amts = append(amts, 10) //1e8)
	}

	/*
		balance, err := TestWallet.GetAddressBalance(recAnp.Address)
		if err != nil {
			t.Fatal(err)
		}
	*/

	//20000e8
	trans, _, err := TestWallet.ConstructTransaction(recs, amts)
	if err != nil {
		t.Fatal(err)
	}

	var total uint64
	var amtsStrs []string
	for _, a := range amts {
		amtsStrs = append(amtsStrs, fmt.Sprintf("%d", a))
	}

	tStruct := TestWallet.Wallet.GetTransactions()[trans]
	total, err = tStruct.TotalInputs()
	if err != nil {
		t.Fatal(err)
	}

	nameComp, err := TestWallet.CheckTransactionAndGetName(recs, amtsStrs, "")
	if err != nil {
		t.Fatal(err)
	} else if trans != nameComp {
		t.Fatal("Names do not match")
	}

	_, err = TestWallet.SendTransaction(trans)
	if err != nil {
		t.Fatal(err)
	}

	// Test string versions
	name, ret, err := TestWallet.ConstructConvertEntryCreditsStrings(recs, amtsStrs)
	if err != nil {
		t.Fatal(err)
	}

	if ret.Total+ret.Fee != total {
		t.Fatal("Total wrong")
	}

	_, err = TestWallet.SendTransaction(name)
	if err != nil {
		t.Fatal(err)
	}

	//_ = balance
	/*if balance != balance+uint64(1e8) {
		t.Error("Balance not changed")
	}*/

	_ = anp
}

// do 8089
func LoadTestWallet(port int) error {
	if TestWallet != nil { // If already instantiated
		return nil
	}

	GUI_DB = MAP
	WALLET_DB = MAP
	TX_DB = MAP

	wal, err := TestHelper.Start(port)
	if err != nil {
		return err
	}

	TestWallet = wal
	return nil
}
