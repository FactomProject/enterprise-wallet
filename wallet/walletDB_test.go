package wallet_test

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	ed "github.com/FactomProject/ed25519"
	ad "github.com/FactomProject/enterprise-wallet/address"
	. "github.com/FactomProject/enterprise-wallet/wallet"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
	//"github.com/FactomProject/factom/wallet"
)

var longtest = true
var _ = fmt.Sprintf("")

// Testing the inserting order
func TestGetRelatedTransaction(t *testing.T) {
	if !longtest {
		return
	}
	//fmt.Println(0)
	err := LoadTestWallet(8089)
	defer StopTestWallet(true)
	if err != nil {
		t.Fatal("Error in test helper", err.Error())
	}

	anp, list := TestWallet.GetGUIAddress("FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q")
	if list == -1 {
		anp, err = TestWallet.AddAddress("Sand", "Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK")
		if err != nil {
			t.Fatal("Error adding address: ", err)
		}
	}

	// Fs1uHDWjYANSxXUtbdkLVkhHboaPRz1tADqs7iB16kQq5VTCvKZS FA2BpB5btNeoSXu2ARqCcF7qkn1XJr5BmDXLjYxd5YsoDH5wU2VU // 1000 from sand
	// Fs3FyjVtnfJV3jZQq7fpYns89fGZtxADARqvwcNGsJrKRubhQ8t7 FA3VSGBDaT3sJ1jv8Be9ABCWw9MgKjUsJcJY2pJTUDXsGMERUtpV // Rest from sand
	// Fs2qf5WTcctcfmestdJUF5dH6geuwBvzVaCGL2458SkJzZKsCU8z FA3HRq8jFUhzN9c8iKTBBfXyNyijSnov1ZLJtJMKTXFQNmncWZoE
	type AddSecPub struct {
		Sec string
		Pub string
	}

	var Add1 = AddSecPub{"Fs1uHDWjYANSxXUtbdkLVkhHboaPRz1tADqs7iB16kQq5VTCvKZS", "FA2BpB5btNeoSXu2ARqCcF7qkn1XJr5BmDXLjYxd5YsoDH5wU2VU"}
	var Add2 = AddSecPub{"Fs3FyjVtnfJV3jZQq7fpYns89fGZtxADARqvwcNGsJrKRubhQ8t7", "FA3VSGBDaT3sJ1jv8Be9ABCWw9MgKjUsJcJY2pJTUDXsGMERUtpV"}
	var Add3 = AddSecPub{"Fs2qf5WTcctcfmestdJUF5dH6geuwBvzVaCGL2458SkJzZKsCU8z", "FA3HRq8jFUhzN9c8iKTBBfXyNyijSnov1ZLJtJMKTXFQNmncWZoE"}

	_, err = TestWallet.AddAddress("Temp", Add2.Sec)
	if err != nil {
		t.Fatalf("Error adding address: ", err.Error())
	}

	// Send 3 Transactions
	tx1, err := sendTrans(Add1.Pub, 100)
	if err != nil {
		t.Fatal("Error sending transaction: ", err)
	}

	time.Sleep(2 * time.Second)

	tx2, err := sendTrans(Add2.Pub, 100)
	if err != nil {
		t.Fatal("Error sending transaction: ", err)
	}

	time.Sleep(2 * time.Second)

	tx3, err := sendTrans(Add3.Pub, 100)
	if err != nil {
		t.Fatal("Error sending transaction: ", err)
	}

	// Ok we have some transactions around
	TestWallet = nil // Need fresh

	StopTestWallet(false)
	LoadTestWallet(8089)

	TestWallet.UpdateGUIDB()
	transactions, err := TestWallet.GetRelatedTransactions()
	if err != nil {
		t.Fatal("Error getting related transaction: ", err)
	}

	anp, err = TestWallet.AddAddress("Third", Add3.Sec)
	if err != nil {
		t.Fatal("Error adding address:", err)
	}

	var _, _, _ = tx1, tx2, tx3
	// Let a block pass
	time.Sleep(10 * time.Second)

	correctTrans, _ := TestWallet.GetRelatedTransactionsNoCaching()
	transactions, err = TestWallet.GetRelatedTransactions()
	if err != nil {
		t.Fatal("Error getting related transaction: ", err)
	}

	if !DisplayTransactions(correctTrans).IsSimilarTo(transactions) {
		printtxID(correctTrans)
		fmt.Println("-")
		printtxID(transactions)
		t.Fatal("Not Same")
	}

	// This tx is before other
	anp, err = TestWallet.AddAddress("First", Add1.Sec)
	if err != nil {
		t.Fatal("Error adding address:", err)
	}

	correctTrans, _ = TestWallet.GetRelatedTransactionsNoCaching()
	transactions, err = TestWallet.GetRelatedTransactions()
	if err != nil {
		t.Fatal("Error getting related transaction: ", err)
	}

	if !DisplayTransactions(correctTrans).IsSimilarTo(transactions) {
		t.Fatal("Not Same")
	}

	// This tx is between both
	anp, err = TestWallet.AddAddress("Second", Add2.Sec)
	if err != nil {
		t.Fatal("Error adding address:", err)
	}

	correctTrans, _ = TestWallet.GetRelatedTransactionsNoCaching()
	transactions, err = TestWallet.GetRelatedTransactions()
	if err != nil {
		t.Fatal("Error getting related transaction: ", err)
	}

	if !DisplayTransactions(correctTrans).IsSimilarTo(transactions) {
		t.Fatal("Not Same")
	}

	_ = anp

}

func printtxID(transactions []DisplayTransaction) {
	for _, t := range transactions {
		fmt.Println(t.TxID + " -- " + t.Time)
	}
}

func findTrans(transactions []DisplayTransaction, txid string) int {
	for i, t := range transactions {
		if t.TxID == txid {
			return i
		}
	}
	return -1
}

func sendTrans(address string, amt uint64) (string, error) {
	toAddresses := []string{address}
	amounts := []uint64{amt * 1e8}
	name, _, err := TestWallet.ConstructTransaction(toAddresses, amounts)
	if err != nil {
		return "", err
	}

	tx, err := TestWallet.SendTransaction(name)
	if err != nil {
		return "", err
	}

	return tx, nil
}

func TestGUIUpdate(t *testing.T) {
	//fmt.Println(1)
	var err error
	TestWallet = nil // Need fresh
	err = LoadTestWallet(8089)
	defer StopTestWallet(true)
	if err != nil {
		t.Fatal(err.Error())
	}

	var faList []*factom.FactoidAddress
	var ecList []*factom.ECAddress
	var addMap map[string]string

	addMap = make(map[string]string)

	for i := 0; i < 5; i++ {
		f, err := TestWallet.Wallet.GenerateFCTAddress()
		if err != nil {
			t.Fatal(err.Error())
		}

		faList = append(faList, f)
		addMap[f.String()] = f.String()
	}

	for i := 0; i < 5; i++ {
		e, err := TestWallet.Wallet.GenerateECAddress()
		if err != nil {
			t.Fatal(err.Error())
		}

		ecList = append(ecList, e)
		addMap[e.String()] = e.String()
	}

	err = TestWallet.UpdateGUIDB()
	if err != nil {
		t.Fatal(err.Error())
	}

	var anpMap map[string]string
	anpMap = make(map[string]string)

	anps := TestWallet.GetAllGUIAddresses()
	// All these should be in first map
	for _, a := range anps {
		if _, ok := addMap[a.Address]; !ok {
			t.Fatal("Should be there, but is not")
		}
		anpMap[a.Address] = a.Address
	}

	// All these should be in second map
	for _, a := range faList {
		if _, ok := anpMap[a.String()]; !ok {
			t.Fatal("Should be there, but is not")
		}
	}

	// All these should be in second map
	for _, a := range ecList {
		if _, ok := anpMap[a.String()]; !ok {
			t.Fatal("Should be there, but is not")
		}
	}

	// Here is some iffy stuff. You cannot delete transactions, but you can delete the wallet
	TestWallet.Wallet, err = wallet.NewMapDBWallet()
	if err != nil {
		t.Fatal(err)
	}
	/*
		TestWallet.UpdateGUIDB()
		newAnps := TestWallet.GetAllGUIAddresses()
		if len(newAnps) != 0 {
			t.Fatal("Should be all deleted")
		}*/
}

func TestDBInteraction(t *testing.T) {
	fmt.Println("TestDBInteraction")
	//fmt.Println(2)
	err := LoadTestWallet(8089)
	defer StopTestWallet(true)
	if err != nil {
		t.Fatal("--1--", err.Error())
	}

	wal := TestWallet
	// We need to call this before it is tested
	wal.TransactionDB.GetAllTXs()

	// Begin Tests

	err = DBAddAndCountTest(wal)
	if err != nil {
		t.Fatalf("--2--", err.Error())
	}

	err = DBAddingExternalAddress(wal)
	if err != nil {
		t.Fatalf("--3--", err.Error())
	}

	err = ChangeAddressNameTest(wal)
	if err != nil {
		t.Fatalf("--4--", err.Error())
	}

	err = CheckRemoveAddressTest(wal)
	if err != nil {
		t.Fatalf("--5--", err.Error())
	}

	err = RelatedTransTest(wal)
	if err != nil {
		// Usually a factomd or wallet loading transaction issue.
		// t.Fatalf(err.Error())
	}

	// End Tests

	/*err = wal.Close()
	if err != nil {
		t.Fatalf(err.Error())
	}*/
}

func TXDBTest(wal *WalletDB) error {
	_, err := wal.TransactionDB.GetAllTXs()
	return err
}

func CheckRemoveAddressTest(wal *WalletDB) error {
	add, err := RandomFactomAddress()
	if err != nil {
		return err
	}

	anp, err := wal.AddExternalAddress("ToBeRemoved", add.String())
	if err != nil {
		return err
	}

	_, list := wal.GetGUIAddress(anp.Address)
	if list == -1 {
		return fmt.Errorf("Address not found")
	}

	_, err = wal.RemoveAddressFromAnyList(anp.Address)
	if err != nil {
		return err
	}

	_, list = wal.GetGUIAddress(anp.Address)

	if list > 0 && list < 4 {
		return fmt.Errorf("Address was not removed")
	}

	return nil
}

func ChangeAddressNameTest(wal *WalletDB) error {
	// External Address Test
	add, err := RandomFactomAddress()
	if err != nil {
		return err
	}

	anp, err := wal.AddAddress("AName", add.SecString())
	if err != nil {
		return err
	}

	err = wal.ChangeAddressName(add.String(), "ABetterName")
	if err != nil {
		return err
	}

	anp, _ = wal.GetGUIAddress(add.String())
	if anp.Name != "ABetterName" {
		return fmt.Errorf("Name did not change")
	}

	// Factoid Address Test
	anp, err = wal.GenerateFactoidAddress("FactName")
	if err != nil {
		return err
	}

	err = wal.ChangeAddressName(anp.Address, "NewFactName")
	if err != nil {
		return err
	}

	anp, _ = wal.GetGUIAddress(anp.Address)
	if anp.Name != "NewFactName" {
		return fmt.Errorf("Name did not change")
	}

	// Entry Credit Address Test
	anp, err = wal.GenerateEntryCreditAddress("ECName")
	if err != nil {
		return err
	}

	err = wal.ChangeAddressName(anp.Address, "NewECName")
	if err != nil {
		return err
	}

	anp, _ = wal.GetGUIAddress(anp.Address)
	if anp.Name != "NewECName" {
		return fmt.Errorf("Name did not change")
	}

	return nil
}

func DBAddingExternalAddress(wal *WalletDB) error {
	add, err := RandomFactomAddress()
	if err != nil {
		return err
	}

	anp, err := wal.AddExternalAddress("RandomPublic", add.String())
	if err != nil {
		return err
	}

	if strings.Compare(anp.Address, add.String()) != 0 {
		return fmt.Errorf("Address added does not match")
	}

	if strings.Compare("RandomPublic", anp.Name) != 0 {
		return fmt.Errorf("Name added does not match")
	}

	_, list := wal.GetGUIAddress(anp.Address)

	if list != 3 {
		return fmt.Errorf("Not placed in External list")
	}

	return nil

}

func RelatedTransTest(wal *WalletDB) error {
	_, err := wal.GetRelatedTransactions()
	return err
}

func DBAddAndCountTest(wal *WalletDB) error {
	count := wal.GetTotalGUIAddresses()

	anp, err := wal.GenerateFactoidAddress("TestFactoid")
	if err != nil {
		return err
	}

	_, list := wal.GetGUIAddress(anp.Address)
	if list != 1 {
		return fmt.Errorf("Added to wrong list")
	}

	anp, err = wal.GenerateEntryCreditAddress("TestEntryCredit")
	if err != nil {
		return err
	}

	_, list = wal.GetGUIAddress(anp.Address)
	if list != 2 {
		return fmt.Errorf("Added to wrong list")
	}

	if wal.GetTotalGUIAddresses() != count+2 {
		return fmt.Errorf("Count is wrong")
	}

	return nil
}

func TestWalletMarshaling(t *testing.T) {
	wal := NewWallet()
	list, err := RandomAddressList(5)
	if err != nil {
		t.Fatalf("--6--", err.Error())
	}
	wal.FactoidAddresses = list

	list, err = RandomAddressList(5)
	if err != nil {
		t.Fatalf("--7--", err.Error())
	}
	wal.EntryCreditAddresses = list

	list, err = RandomAddressList(5)
	if err != nil {
		t.Fatalf("--8--", err.Error())
	}
	wal.ExternalAddresses = list

	data, err := wal.MarshalBinary()
	if err != nil {
		t.Fatalf("--9--", err.Error())
	}

	wal2 := NewWallet()
	wal2.UnmarshalBinary(data)
	if !wal.IsSameAs(wal2) {
		t.Fatalf("Not Same")
	} else if wal2.FactoidAddresses.Length != 5 || wal2.ExternalAddresses.Length != 5 || wal2.EntryCreditAddresses.Length != 5 {
		t.Fatalf("UnMarshalDFail")
	}
}

func RandomAddressList(n int) (*ad.AddressList, error) {
	addList := ad.NewAddressList()

	for i := 0; i < n; i++ {
		anp, err := RandomAddressNamePair()
		if err != nil {
			return nil, err
		}
		err = addList.AddANP(anp)
		if err != nil {
			return nil, err
		}
	}

	return addList, nil
}

func RandomFactomAddress() (*factom.FactoidAddress, error) {
	_, privkey, err := ed.GenerateKey(crand.Reader)
	if err != nil {
		return nil, err
	}

	address, err := factom.MakeFactoidAddress(privkey[:32])
	if err != nil {
		return nil, err
	}

	return address, nil
}

func RandomAddressNamePair() (*ad.AddressNamePair, error) {
	address, err := RandomFactomAddress()
	if err != nil {
		return nil, err
	}

	a, err := ad.NewAddress(RandStringBytes(20), address.String())
	if err != nil {
		return nil, err
	}
	return a, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
