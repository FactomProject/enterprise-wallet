package address_test

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	ed "github.com/FactomProject/ed25519"
	. "github.com/FactomProject/enterprise-wallet/address"
	"github.com/FactomProject/factom"
)

var _ = strings.Compare("", "")
var _ = fmt.Sprintf("")

func TestAddressNamePairMarshal(t *testing.T) {
	a, err := NewSeededAddress("Factoid1ss", "FA27kaVcH76hDsLmZuSq2yad6zrmUDUm6KCHq6nibEZiKbBSLQ8C")
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

	c := new(AddressNamePair)
	err = c.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("Failed UnMarshal Binary: %s", err.Error())
	}

	if !b.IsSameAs(a) {
		t.Fatalf("Failed: Not same")
	}
	if !c.IsSameAs(b) {
		t.Fatalf("Failed: Not same")
	}
}

func TestNewAddressFails(t *testing.T) {
	add, err := RandomAddress()
	if err != nil {
		t.Fatalf("Failed random address: %s", err.Error())
	}

	_, err = NewSeededAddress(
		"12345678901234567890111", // Over 20 characters
		add)

	if err == nil {
		t.Fatalf("Failed new address accepts too long of a name")
	}

	add = "FA38rrBLKrbt7qCLFYUiF8yg8zYKL5GMkWysJV8H45Rux5aBY9cf" // Invalid

	_, err = NewAddress(
		"1234567890123456789", // Over 20 characters
		add)
	if err == nil {
		t.Fatalf("Failed new address accepts bad address")
	}

	anp, err := RandomAddressNamePair()
	if err != nil {
		t.Fatalf("Failed random anp: %s", err.Error())
	}

	err = anp.ChangeName("123456789012345678901")
	if err == nil {
		t.Fatalf("Failed change address name accepts too long of a name")
	}

	anp2, _ := RandomAddressNamePair()
	na, _ := RandomAddress()
	anp2.Address = na
	anp.ChangeName("newname")
	anp2.ChangeName("Newname")

	if anp.IsSameAs(anp2) {
		t.Fatalf("Failed not same")
	}

}

func TestAddressList(t *testing.T) {
	addList := NewAddressList()

	for i := 0; i < 10; i++ {
		anp, err := RandomAddressNamePair()
		if err != nil {
			t.Fatalf(err.Error())
		}
		_, err = addList.AddSeeded(anp.Name, anp.Address)
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

	addList.List[0].Name = "REMOVEME"
	x := addList.List[0].Address
	err = addList.Remove(x)
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, i := addList.Get(x)
	if i != -1 {
		t.Fatalf("Should have been removed")
	}

	if addList2.IsSameAs(addList) {
		t.Fatalf("Failed: Not same")
	}

	err = addList2.Remove(addList2.List[0].Address)
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

func TestAddressListFunctions(t *testing.T) {
	addList := NewAddressList()
	anpSave, err := RandomAddressNamePair()
	if err != nil {
		t.Fatalf(err.Error())
	}

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

	_, err = addList.Add(anpSave.Name, anpSave.Address)
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, i := addList.Get(anpSave.Address)

	addList.List[i].ChangeName("ThisIsABetterName")

	newAnp, _ := addList.Get(anpSave.Address)
	if newAnp.Name != "ThisIsABetterName" {
		t.Fatalf("name did not change")
	}

	anp, err := RandomAddressNamePair()
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = addList.AddANP(anp)
	if err != nil {
		t.Fatalf("Cannot add new anp %s\n", err.Error())
	}

	err = addList.AddANP(anp)
	if err == nil {
		t.Fatalf("Added duplicate anp, should not\n")
	}

	anp.Address = "FA38rrBLKrbt7qCLFYUiF8yg8zYKL5GMkWysJV8H45Rux5aBY9cf"
	err = addList.AddANP(anp)
	if err == nil {
		t.Fatalf("Added invalid address anp\n")
	}

	anp.Address, err = RandomAddress()
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = addList.AddANP(anp)
	if err != nil {
		t.Fatalf("Did not add valid address anp\n")
	}

	_, err = addList.Add("", anp.Address)
	if err == nil {
		t.Fatalf("Added nil name\n")
	}

	anp.Address = "FA38rrBLKrbt7qCLFYUiF8yg8zYKL5GMkWysJV8H45Rux5aBY9cf"
	_, err = addList.Add("Good", anp.Address)
	if err == nil {
		t.Fatalf("Added invalid address anp\n")
	}

	err = addList.Remove(anp.Address)
	if err == nil {
		t.Fatalf("Removed non-existent anp\n")
	}

	anp2, err := RandomAddressNamePair()
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = addList.Add(anp2.Name, anp2.Address)
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = addList.Add(anp2.Name, anp2.Address)
	if err == nil {
		t.Fatalf("Added duplicate address")
	}

	err = addList.Remove(anp2.Address)
	if err != nil {
		t.Fatalf(err.Error())
	}

	anp.Address = "FA38rrBLKrbt7qCLFYUiF8yg8zYKL5GMkWysJV8H45Rux5aBY9cf"
	err = addList.Remove(anp.Address)
	if err == nil {
		t.Fatalf("Removed bada ddress")
	}

	//
	a, _ := RandomAddressNamePair()
	aL := NewAddressList()
	bL := NewAddressList()

	aL.Add("AName", a.Address)
	bL.Add("BName", a.Address)
	if aL.IsSameAs(bL) {
		t.Fatalf("Not same\n")
	}

}

func RandomAddressNamePair() (*AddressNamePair, error) {
	add, err := RandomAddress()
	if err != nil {
		return nil, err
	}

	a, err := NewAddress(RandStringBytes(20), add)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func RandomAddress() (string, error) {
	_, privkey, err := ed.GenerateKey(crand.Reader)
	if err != nil {
		return "", err
	}

	address, err := factom.MakeFactoidAddress(privkey[:32])
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
