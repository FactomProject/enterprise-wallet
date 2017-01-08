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

func TestGUIWalletFromWallet(t *testing.T) {
	w := NewWallet()

	// Not valid addresses
	_, err := w.AddAddress("Test", "NotAValid", 1)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	_, err = w.AddAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s6Z2cX6cakdQ6U53uSLRoPLR", 1)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	_, err = w.AddAddress("Test", "EC32x8uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 1)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	_, err = w.AddAddress("Test", "EC2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 1)
	if err == nil {
		t.Fatal("Accepted an invalid address")
	}

	// Valid wrong list
	_, err = w.AddAddress("Test", "EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 1)
	if err == nil {
		t.Fatal("Accepted a valid address in wrong list")
	}

	_, err = w.AddAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 2)
	if err == nil {
		t.Fatal("Accepted a valid address in wrong list")
	}

	// Valid address, invalid list
	_, err = w.AddAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 0)
	if err == nil {
		t.Fatal("Accepted an invalid list")
	}

	_, err = w.AddAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 4)
	if err == nil {
		t.Fatal("Accepted an invalid list")
	}

	_, err = w.AddAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", -1)
	if err == nil {
		t.Fatal("Accepted an invalid list")
	}

	// Valid
	_, err = w.AddAddress("Test", "FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR", 1)
	if err != nil {
		t.Fatal("Rejected a valid address")
	}

	_, err = w.AddAddress("Test", "EC3FmWu7iX85r6UvTaqBEZgNNGAmNE1Vd2ZXRGaxHr1g8jRcS6TQ", 2)
	if err != nil {
		t.Fatal("Rejected a valid address")
	}

	_, err = w.AddAddress("Test", "EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P", 3)
	if err != nil {
		t.Fatal("Rejected a valid address")
	}

	_, err = w.AddAddress("Test", "FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb", 3)
	if err != nil {
		t.Fatal("Rejected a valid address")
	}
}

func TestGUIWallet(t *testing.T) {
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

}

// Valid
// FA2SDU3UhBwrBR2q7jbFAbxnqUW6s5Z2cX6cakdQ6U53uSLRoPLR
// FA39udanfmkZXZxPUjMWqmXvdUNKSN9D3UCTnNsJX9B4n7dadCUb
// EC32x9uN4xMEMQbw66oob2de94z3b1JWhn23E9srgG3aCzhCCa3P
// EC3FmWu7iX85r6UvTaqBEZgNNGAmNE1Vd2ZXRGaxHr1g8jRcS6TQ
