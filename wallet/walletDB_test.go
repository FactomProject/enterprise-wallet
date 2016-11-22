package wallet_test

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	ad "github.com/FactomProject/M2WalletGUI/address"
	. "github.com/FactomProject/M2WalletGUI/wallet"
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/factom"
	//"github.com/FactomProject/factom/wallet"
)

var _ = fmt.Sprintf("")

func TestDBInteraction(t *testing.T) {
	err := LoadTestWallet()
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
