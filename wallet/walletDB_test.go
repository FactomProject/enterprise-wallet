package wallet_test

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	ad "github.com/FactomProject/M2GUIWallet/address"
	. "github.com/FactomProject/M2GUIWallet/wallet"
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
	//"github.com/FactomProject/factom/wallet"
)

var _ = fmt.Sprintf("")

func TestGUIUpdate(t *testing.T) {
	var err error
	TestWallet = nil // Need fresh
	LoadTestWallet(8079)

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

	TestWallet.UpdateGUIDB()
	newAnps := TestWallet.GetAllGUIAddresses()
	if len(newAnps) != 0 {
		t.Fatal("Should be all deleted")
	}
}

func TestDBInteraction(t *testing.T) {
	err := LoadTestWallet(8089)
	if err != nil {
		t.Fatalf(err.Error())
	}

	wal := TestWallet

	// Begin Tests

	err = DBAddAndCountTest(wal)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = DBAddingExternalAddress(wal)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = RelatedTransTest(wal)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = ChangeAddressNameTest(wal)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = CheckRemoveAddressTest(wal)
	if err != nil {
		t.Fatalf(err.Error())
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

	anp, err := wal.AddAddress("ToBeRemoved", add.SecString())
	if err != nil {
		return err
	}

	_, list := wal.GetGUIAddress(anp.Address)
	if list == -1 {
		return fmt.Errorf("Address not found")
	}

	_, err = wal.RemoveAddress(anp.Address)
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

	anp, err := wal.AddAddress("RandomSecret", add.SecString())
	if err != nil {
		return err
	}

	if strings.Compare(anp.Address, add.String()) != 0 {
		return fmt.Errorf("Address added does not match")
	}

	if strings.Compare("RandomSecret", anp.Name) != 0 {
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
		t.Fatalf(err.Error())
	}
	wal.FactoidAddresses = list

	list, err = RandomAddressList(5)
	if err != nil {
		t.Fatalf(err.Error())
	}
	wal.EntryCreditAddresses = list

	list, err = RandomAddressList(5)
	if err != nil {
		t.Fatalf(err.Error())
	}
	wal.ExternalAddresses = list

	data, err := wal.MarshalBinary()
	if err != nil {
		t.Fatalf(err.Error())
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
