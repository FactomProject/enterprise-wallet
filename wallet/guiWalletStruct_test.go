package wallet_test

import (
	//crand "crypto/rand"
	//"fmt"
	//"math/rand"
	//"strings"
	"testing"

	//ad "github.com/FactomProject/enterprise-wallet/address"
	. "github.com/FactomProject/enterprise-wallet/wallet"
	//ed "github.com/FactomProject/ed25519"
	//"github.com/FactomProject/factom"
	//"github.com/FactomProject/factom/wallet"
)

func TestGUIAddAddress(t *testing.T) {
	w := NewWallet()

	// Not valid addresses
	_, err := w.AddSeededAddress("Test", "NotAValid", 1)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	_, err = w.AddSeededAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s6Z2cX6cakdQ6U53uSLRoPLR", 1)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	_, err = w.AddSeededAddress("Test", "EC32x8uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 1)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	_, err = w.AddSeededAddress("Test", "EC2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 1)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	_, err = w.AddSeededAddress("Test", "NA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 3)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	// Valid wrong list
	_, err = w.AddSeededAddress("Test", "EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 1)
	if err == nil {
		t.Fatal("Accepted a valid address in wrong list")
	}

	_, err = w.AddSeededAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 2)
	if err == nil {
		t.Fatal("Accepted a valid address in wrong list")
	}

	// Valid address, invalid list
	_, err = w.AddSeededAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 0)
	if err == nil {
		t.Fatal("Accepted an invalid list")
	}

	_, err = w.AddSeededAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 4)
	if err == nil {
		t.Fatal("Accepted an invalid list")
	}

	_, err = w.AddSeededAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", -1)
	if err == nil {
		t.Fatal("Accepted an invalid list")
	}

	// Valid
	_, err = w.AddSeededAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 1)
	if err != nil {
		t.Fatal("Rejected a valid address")
	}

	_, err = w.AddSeededAddress("Test", "EC3FmWu7iX85r6UvTaqBEZgNNGAmNE1Vd2ZXRGaxHr1g8jRcS6TQ", 2)
	if err != nil {
		t.Fatal("Rejected a valid address")
	}

	_, err = w.AddSeededAddress("Test", "EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 3)
	if err != nil {
		t.Fatal("Rejected a valid address")
	}

	_, err = w.AddSeededAddress("Test", "FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 3)
	if err != nil {
		t.Fatal("Rejected a valid address")
	}

	// Adding again is invalid
	_, err = w.AddSeededAddress("Test", "FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 3)
	if err == nil {
		t.Fatal("Added the same address twice, this is not supposed to be allowed")
	}

	w.AddBalancesToAddresses()
}

func TestGetAddress(t *testing.T) {
	gw := NewWallet()
	_, err := gw.AddAddress("1", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = gw.AddAddress("1", "FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 3)
	if err != nil {
		t.Fatal(err)
	}
	_, err = gw.AddAddress("1", "EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = gw.AddAddress("1", "EC3FmWu7iX85r6UvTaqBEZgNNGAmNE1Vd2ZXRGaxHr1g8jRcS6TQ", 2)
	if err != nil {
		t.Fatal(err)
	}

	list := gw.GetAllAddresses()
	if len(list) != 4 {
		t.Fatal("List wrong length")
	}

	count := gw.GetTotalAddressCount()
	if count != 4 {
		t.Fatal("List wrong length")
	}

	list = gw.GetAllAddressesFromList(1)
	if len(list) != 1 {
		t.Fatal("List wrong length")
	}
	list = gw.GetAllAddressesFromList(2)
	if len(list) != 2 {
		t.Fatal("List wrong length")
	}
	list = gw.GetAllAddressesFromList(3)
	if len(list) != 1 {
		t.Fatal("List wrong length")
	}

	data, err := gw.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	gw2 := NewWallet()
	err = gw2.UnmarshalBinary(data)
	if err != nil {
		t.Fatal(err)
	}

	if !gw2.IsSameAs(gw) {
		t.Fatal("Not same, but are")
	}

	// All same as g2
	gw2Clones := make([]*WalletStruct, 3)
	for i := 0; i < 3; i++ {
		gw2Clones[i] = new(WalletStruct)
		gw2Clones[i].UnmarshalBinaryData(data)
	}

	gw2Clones[0].FactoidAddresses.AddSeeded("RandomAdditonal", "Fs1rMCawoAk26xMa7WMxnbVNW8w69wQ7kRYYTgCrosFTqAJ19wwZ")
	gw2Clones[1].EntryCreditAddresses.AddSeeded("RandomAdd", "EC3JFMMpSpDEZFf7hBeSrcx25s6jkkoCV1F654J1uruBxNZRKCvF")
	gw2Clones[2].ExternalAddresses.AddSeeded("RandomAdd", "EC3JFMMpSpDEZFf7hBeSrcx25s6jkkoCV1F654J1uruBxNZRKCvF")

	for i := range gw2Clones {
		if gw2Clones[i].IsSameAs(gw2) {
			t.Error("Failed: Not the same")
		}
	}
}

func TestChangeName(t *testing.T) {
	gw := NewWallet()
	var err error
	tooLong := "WayTooLongOfANameWayTooLongOfANameWayTooLongOfANameWayTooLongOfANameWayTooLongOfAName"

	err = gw.ChangeAddressName("FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", "2")
	if err == nil {
		t.Fatal("Address does not exist yet")
	}

	// List 1
	_, err = gw.AddAddress("1", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 1)
	if err != nil {
		t.Fatal(err)
	}

	err = gw.ChangeAddressName("FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", "2")
	if err != nil {
		t.Fatal(err)
	}

	anp, _, _ := gw.GetAddress("FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR")
	if anp.Name != "2" {
		t.Error("Name did not change")
	}

	err = gw.ChangeAddressName("FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", tooLong)
	if err == nil {
		t.Fatal("Name should be rejected")
	}

	// List 3
	_, err = gw.AddAddress("1", "EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 2)
	if err != nil {
		t.Fatal(err)
	}

	err = gw.ChangeAddressName("EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", "2")
	if err != nil {
		t.Fatal(err)
	}

	anp, _, _ = gw.GetAddress("EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P")
	if anp.Name != "2" {
		t.Error("Name did not change")
	}

	err = gw.ChangeAddressName("EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", tooLong)
	if err == nil {
		t.Fatal("Name should be rejected")
	}

	// List 2
	_, err = gw.AddAddress("1", "FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 3)
	if err != nil {
		t.Fatal(err)
	}

	err = gw.ChangeAddressName("FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", "2")
	if err != nil {
		t.Fatal(err)
	}

	anp, _, _ = gw.GetAddress("FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb")
	if anp.Name != "2" {
		t.Error("Name did not change")
	}

	err = gw.ChangeAddressName("FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", tooLong)
	if err == nil {
		t.Fatal("Name should be rejected")
	}
}

func TestRemove(t *testing.T) {
	gw := NewWallet()
	var err error
	_, err = gw.RemoveAddress("FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 0)
	if err == nil {
		t.Error("Invalid list")
	}

	_, err = gw.RemoveAddress("FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 1)
	if err == nil {
		t.Error("Address does not exist, should return an error")
	}

	_, err = gw.RemoveAddress("EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 2)
	if err == nil {
		t.Error("Address does not exist, should return an error")
	}

	_, err = gw.RemoveAddress("FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 3)
	if err == nil {
		t.Error("Address does not exist, should return an error")
	}

	// Good removes
	gw.AddSeededAddress("1", "FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 1)
	gw.AddSeededAddress("1", "EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 2)
	gw.AddAddress("1", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 3)

	_, err = gw.RemoveAddress("FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 1)
	if err != nil {
		t.Error(err)
	}

	_, err = gw.RemoveAddress("EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 2)
	if err != nil {
		t.Error(err)
	}

	_, err = gw.RemoveAddress("FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 3)
	if err != nil {
		t.Error(err)
	}
}

// Valid
// FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR
// FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb
// EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P
// EC3FmWu7iX85r6UvTaqBEZgNNGAmNE1Vd2ZXRGaxHr1g8jRcS6TQ
