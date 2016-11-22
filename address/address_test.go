package address_test

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"testing"

	. "github.com/FactomProject/M2WalletGUI/address"
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/factom"
)

var _ = fmt.Sprintf("")

func TestAddressNamePairMarshal(t *testing.T) {
	a, err := NewAddress("Factoid1ss", "FA27kaVcH76hDsLmZuSq2yad6zrmUDUm6KCHq6nibEZiKbBSLQ8C")
	if err != nil {
		t.Fatalf("Error Creating")
	}

	data, err := a.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed Marshal Binary: %s", err.Error())
	}

	b := new(AddressNamePair)
	_, err = b.UnmarshalBinaryData(data)
	if err != nil {
		t.Fatalf("Failed UnMarshal Binary: %s", err.Error())
	}

	if !b.IsSameAs(a) {
		t.Fatalf("Failed: Not same")
	}
}

func TestAddressList(t *testing.T) {
	addList := NewAddressList()

	for i := 0; i < 10; i++ {
		anp, err := RandomAddressNamePair()
		if err != nil {
			t.Fatalf(err.Error())
		}
		_, err = addList.Add(anp.Name, anp.Address)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	// Marshal/Umarshal, add/remove
	data, err := addList.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed Marshal Binary: %s", err.Error())
	}

	addList2 := NewAddressList()
	addList2.UnmarshalBinary(data)

	if addList.Length != 10 || addList2.Length != 10 {
		t.Fatalf("Failed: Bad length")
	}

	if !addList.IsSameAs(addList2) {
		t.Fatalf("Failed: Not same")
	}

	err = addList.Remove(addList.List[0])
	if err != nil {
		t.Fatalf(err.Error())
	}

	if addList2.IsSameAs(addList) {
		t.Fatalf("Failed: Not same")
	}

	err = addList2.Remove(addList2.List[0])
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !addList.IsSameAs(addList2) {
		t.Fatalf("Failed: Not same")
	}

	// Check 2 back to back
	data = append(data, data[:]...)
	addList3 := NewAddressList()
	data, err = addList3.UnmarshalBinaryData(data)
	if err != nil {
		t.Fatalf(err.Error())
	}
	addList4 := NewAddressList()
	data, err = addList4.UnmarshalBinaryData(data)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !addList3.IsSameAs(addList4) || addList3.Length != 10 {
		t.Fatalf("Failed: Not same")
	}

	// Test nil list
	addList5 := NewAddressList()
	data, err = addList5.MarshalBinary()
	if err != nil {
		t.Fatalf(err.Error())
	}

	addList6 := NewAddressList()
	err = addList6.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !addList6.IsSameAs(addList5) || addList5.Length != 0 {
		t.Fatalf("Failed: Not same")
	}
}

func RandomAddressNamePair() (*AddressNamePair, error) {
	_, privkey, err := ed.GenerateKey(crand.Reader)
	if err != nil {
		return nil, err
	}

	address, err := factom.MakeFactoidAddress(privkey[:32])
	if err != nil {
		return nil, err
	}

	a, err := NewAddress(RandStringBytes(20), address.String())
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
