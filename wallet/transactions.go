package wallet

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"

	"github.com/FactomProject/factom"
	//"github.com/FactomProject/factom/wallet"
)

var _ = fmt.Sprintf("")
var _ = factom.AddressLength

type AddressBalancePair struct {
	Address string
	Balance uint64
}

type AddressBalancePairs []AddressBalancePair

func (slice AddressBalancePairs) Len() int {
	return len(slice)
}

func (slice AddressBalancePairs) Less(i int, j int) bool {
	return slice[i].Balance < slice[j].Balance
}

func (slice AddressBalancePairs) Swap(i int, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice AddressBalancePairs) Index(i int) AddressBalancePair {
	return slice[i]
}

func (wal *WalletDB) ConstructSendFactoidsStrings(toAddresses []string, amounts []string) (string, error) {
	var amts []uint64
	for _, a := range amounts {
		amt64, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return "", err
		}
		amts = append(amts, uint64(amt64*1e8))
	}

	return wal.ConstructSendFactoids(toAddresses, amts)
}

// Constructs factoid transaction
// Transaction name is hash of all the addresses. More than 1 transaction to
// an address(es) should not be open, but combined.
// Parameters:
//		toAddresses = list of output addresses
//		amounts = list of amounts to each output, indicies must match
// Returns:
//		Transaction Name, error
func (wal *WalletDB) ConstructSendFactoids(toAddresses []string, amounts []uint64) (string, error) {
	if len(toAddresses) != len(amounts) {
		return "", fmt.Errorf("Lengths of address to amount does not match")
	} else if len(toAddresses) == 0 {
		return "", fmt.Errorf("No recepient given")
	}

	// Add outputs, find total being sent
	trans := hashStringList(toAddresses)
	trans = trans[:32]
	err := wal.Wallet.NewTransaction(trans)
	if err != nil {
		return "", err
	}

	var total uint64 = 0
	for i, address := range toAddresses {
		wal.Wallet.AddOutput(trans, address, amounts[i])
		total += amounts[i]
	}

	// Decide what addresses to pay with
	// Pay with largest first
	faAddresses, err := wal.Wallet.GetAllFCTAddresses()
	if err != nil {
		return "", err
	}

	var balances []AddressBalancePair

	for _, address := range faAddresses {
		addr := address.String()
		balance, err := factom.GetFactoidBalance(addr)
		if err != nil {
			return "", err
		}
		balances = append(balances, AddressBalancePair{addr, uint64(balance)})
	}

	balList := AddressBalancePairs(balances)
	sort.Sort(sort.Reverse(balList))
	list := []AddressBalancePair(balList)

	totalLeft := total
	var i int
	for i = 0; totalLeft > 0; {
		if i >= len(list) {
			return "", fmt.Errorf("Not enough factoids to cover the transaction")
		}
		if list[i].Balance > totalLeft {
			wal.Wallet.AddInput(trans, list[i].Address, totalLeft)
			list[i].Balance -= totalLeft
			totalLeft = 0
		} else {
			if list[i].Balance > 0 {
				wal.Wallet.AddInput(trans, list[i].Address, list[i].Balance)
				totalLeft -= list[i].Balance
				list[i].Balance = 0
			}
			i++
		}
	}

	transStruct := wal.Wallet.GetTransactions()[trans]
	if transStruct == nil {
		return "", fmt.Errorf("Transaction not found")
	}

	rate, err := factom.GetRate()
	if err != nil {
		return "", err
	}

	fee, err := transStruct.CalculateFee(rate)
	if err != nil {
		return "", err
	}

	for list[i].Balance < fee {
		i++
		if i >= len(list) {
			return "", fmt.Errorf("Not enough factoids to cover the transaction")
		}
	}

	err = wal.Wallet.AddFee(trans, list[i].Address, rate)
	if err != nil {
		return "", err
	}

	err = wal.Wallet.SignTransaction(trans)
	if err != nil {
		return "", err
	}

	return trans, nil
}

// Constructs Convert Entry Credits
// Transaction name is hash of all the addresses. This case only 1 address
// Parameters:
//		toAddresse = entry credit address
//		amount = amount of EC to send
// Returns:
//		Transaction Name, error
func (wal *WalletDB) ConstructConvertToEC(toAddresses []string, amounts []uint64) (string, error) {
	if len(toAddresses) != len(amounts) {
		return "", fmt.Errorf("Lengths of address to amount does not match")
	} else if len(toAddresses) == 0 {
		return "", fmt.Errorf("No recepient given")
	}

	// Add outputs, find total being sent
	trans := hashStringList(toAddresses)
	trans = trans[:32]
	err := wal.Wallet.NewTransaction(trans)
	if err != nil {
		return "", err
	}

	// Cost in Factoids
	rate, err := factom.GetRate()
	if err != nil {
		return "", err
	}

	var total uint64 = 0
	for i, address := range toAddresses {
		amt := rate * amounts[i]
		wal.Wallet.AddECOutput(trans, address, amt)
		total += amt
	}

	// Decide what addresses to pay with
	// Pay with largest first
	faAddresses, err := wal.Wallet.GetAllFCTAddresses()
	if err != nil {
		return "", err
	}

	var balances []AddressBalancePair

	for _, address := range faAddresses {
		addr := address.String()
		balance, err := factom.GetFactoidBalance(addr)
		if err != nil {
			return "", err
		}
		balances = append(balances, AddressBalancePair{addr, uint64(balance)})
	}

	balList := AddressBalancePairs(balances)
	sort.Sort(sort.Reverse(balList))
	list := []AddressBalancePair(balList)

	var totalLeft uint64 = total
	var i int
	for i = 0; totalLeft > 0; {
		if i >= len(balList) {
			return "", fmt.Errorf("Not enough factoids to cover the transaction")
		}
		if balList[i].Balance > totalLeft {
			wal.Wallet.AddInput(trans, balList[i].Address, totalLeft)
			balList[i].Balance -= totalLeft
			totalLeft = 0
		} else if balList[i].Balance < totalLeft {
			if balList[i].Balance > 0 {
				wal.Wallet.AddInput(trans, balList[i].Address, balList[i].Balance)
				totalLeft -= balList[i].Balance
				balList[i].Balance = 0
			}
			i++
		}
	}

	transStruct := wal.Wallet.GetTransactions()[trans]

	fee, err := transStruct.CalculateFee(rate)
	if err != nil {
		return "", err
	}

	for i < len(balList) {
		if list[i].Balance < fee {
			i++
		} else {
			break
		}
	}

	err = wal.Wallet.AddFee(trans, list[i].Address, rate)
	if err != nil {
		return "", err
	}

	err = wal.Wallet.SignTransaction(trans)
	if err != nil {
		return "", err
	}

	return trans, nil
}

func (wal *WalletDB) GetAddressBalance(address string) (uint64, error) {
	bal, err := factom.GetFactoidBalance(address)
	return uint64(bal), err
}

func (wal *WalletDB) SendTransaction(trans string) (string, error) {
	return factom.SendTransaction(trans)
}

func hashStringList(list []string) string {
	buf := new(bytes.Buffer)
	for _, item := range list {
		data := sha256.Sum256([]byte(item))
		buf.Write(data[:])
	}

	trans := sha256.Sum256(buf.Next(buf.Len()))
	return hex.EncodeToString(trans[:])
}
