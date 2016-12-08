package wallet

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
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

// Doublechecks the transaction is the same (with amounts and addresses)
// This is to confirm an already constructed transaction
func (wal *WalletDB) CheckTransactionAndGetName(toAddresses []string, amounts []string) (string, error) {
	name := hashStringList(toAddresses)
	name = name[:32] // name of transaction

	trans := wal.Wallet.GetTransactions()
	t, ok := trans[name]
	if ok || t != nil {
		var outs []interfaces.ITransAddress
		faOuts := t.GetOutputs()
		ecOuts := t.GetECOutputs()

		for _, f := range faOuts {
			outs = append(outs, f)
		}
		for _, e := range ecOuts {
			outs = append(outs, e)
		}

		if len(outs) != len(amounts) {
			return name, fmt.Errorf("A change in the amount of outputs has been detected")
		}
		for i, o := range outs {
			amt, err := strconv.Atoi(amounts[i])
			if err != nil {
				return name, fmt.Errorf("Amount was not able to be converted into a number")
			}

			compAddr := ""
			if toAddresses[i][:2] == "FA" {
				compAddr = primitives.ConvertFctAddressToUserStr(o.GetAddress())
				amt = amt * 1e8
				if o.GetAmount() != uint64(amt) {
					return name, fmt.Errorf("A change in the amount of an output has been detected")
				}
			} else {
				compAddr = primitives.ConvertECAddressToUserStr(o.GetAddress())
				// Compare amt, but rate changes
			}

			if compAddr != toAddresses[i] {
				return name, fmt.Errorf("A change in the address of an output has been detected")
			}
		}
	} else {
		return name, fmt.Errorf("No transaction found that matches the output addresses")
	}
	return name, nil
}

type ReturnTransStruct struct {
	Total uint64 `json:"Total"`
	Fee   uint64 `json:"Fee"`
}

func (wal *WalletDB) ConstructSendFactoidsStrings(toAddresses []string, amounts []string) (string, *ReturnTransStruct, error) {
	var amts []uint64
	for _, a := range amounts {
		amt64, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return "", nil, err
		}
		amts = append(amts, uint64(amt64*1e8))
	}

	return wal.ConstructSendFactoids(toAddresses, amts)
}

func (wal *WalletDB) ConstructConvertEntryCreditsStrings(toAddresses []string, amounts []string) (string, *ReturnTransStruct, error) {
	var amts []uint64
	for _, a := range amounts {
		amt64, err := strconv.ParseUint(a, 10, 64)
		if err != nil {
			return "", nil, err
		}
		amts = append(amts, amt64)
	}

	return wal.ConstructConvertToEC(toAddresses, amts)
}

func (wal *WalletDB) DeleteTransaction(trans string) error {
	return wal.Wallet.DeleteTransaction(trans)
}

// Constructs factoid transaction
// Transaction name is hash of all the addresses. More than 1 transaction to
// an address(es) should not be open, but combined.
// Parameters:
//		toAddresses = list of output addresses
//		amounts = list of amounts to each output, indicies must match
// Returns:
//		Transaction Name, error
func (wal *WalletDB) ConstructSendFactoids(toAddresses []string, amounts []uint64) (string, *ReturnTransStruct, error) {
	if len(toAddresses) != len(amounts) {
		return "", nil, fmt.Errorf("Lengths of address to amount does not match")
	} else if len(toAddresses) == 0 {
		return "", nil, fmt.Errorf("No recepient given")
	}

	// Add outputs, find total being sent
	trans := hashStringList(toAddresses)
	trans = trans[:32] // Name of transaction

	transMap := wal.Wallet.GetTransactions()
	if t, _ := transMap[trans]; t != nil {
		wal.DeleteTransaction(trans)
	}

	err := wal.Wallet.NewTransaction(trans)
	if err != nil {
		return trans, nil, err
	}

	var total uint64 = 0
	for i, address := range toAddresses {
		if !wal.IsValidAddress(address) || address[:2] != "FA" {
			return trans, nil, fmt.Errorf("Invalid address given")
		}
		wal.Wallet.AddOutput(trans, address, amounts[i])
		total += amounts[i]
	}

	// Decide what addresses to pay with
	// Pay with largest first
	faAddresses, err := wal.Wallet.GetAllFCTAddresses()
	if err != nil {
		return trans, nil, err
	}

	var list []AddressBalancePair

	for _, address := range faAddresses {
		addr := address.String()
		balance, err := factom.GetFactoidBalance(addr)
		if err != nil {
			return trans, nil, err
		}
		list = append(list, AddressBalancePair{addr, uint64(balance)})
	}

	// Sort to get largest balances first
	sort.Sort(sort.Reverse(AddressBalancePairs(list)))

	totalLeft := total
	var i int
	// While factoids still needed to cover transaction, go through addresses
	for i = 0; totalLeft > 0; {
		if i >= len(list) {
			return trans, nil, fmt.Errorf("Not enough factoids to cover the transaction")
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
		return trans, nil, fmt.Errorf("Transaction not found")
	}

	rate, err := factom.GetRate()
	if err != nil {
		return trans, nil, err
	}

	fee, err := transStruct.CalculateFee(rate)
	if err != nil {
		return trans, nil, err
	}

	// The last addresse we used to pay, we need to check if it can cover the fee
	if list[i].Balance < fee { // If it cannot, lets find one that can
		i, err = checkForAddressForFee(list, transStruct, i, rate)
		if i == -1 || err != nil { // We don't have an address that can pay for the fee.
			return trans, nil, fmt.Errorf("Not enough factoids to cover the transaction")
		} else {
			wal.Wallet.AddInput(trans, list[i].Address, 0)
		}
	}

	err = wal.Wallet.AddFee(trans, list[i].Address, rate)
	if err != nil {
		return trans, nil, err
	}

	err = wal.Wallet.SignTransaction(trans)
	if err != nil {
		return trans, nil, err
	}

	r := new(ReturnTransStruct)
	r.Total = total
	r.Fee = fee

	return trans, r, nil
}

// Constructs Convert Entry Credits
// Transaction name is hash of all the addresses. This case only 1 address
// Parameters:
//		toAddresse = entry credit address
//		amount = amount of EC to send
// Returns:
//		Transaction Name, error
func (wal *WalletDB) ConstructConvertToEC(toAddresses []string, amounts []uint64) (string, *ReturnTransStruct, error) {
	if len(toAddresses) != len(amounts) {
		return "", nil, fmt.Errorf("Lengths of address to amount does not match")
	} else if len(toAddresses) == 0 {
		return "", nil, fmt.Errorf("No recepient given")
	}

	// Add outputs, find total being sent
	trans := hashStringList(toAddresses)
	trans = trans[:32] // Name of transaction

	transMap := wal.Wallet.GetTransactions()
	if t, _ := transMap[trans]; t != nil {
		wal.DeleteTransaction(trans)
	}

	err := wal.Wallet.NewTransaction(trans)
	if err != nil {
		return trans, nil, err
	}

	// Cost in Factoids
	rate, err := factom.GetRate()
	if err != nil {
		return trans, nil, err
	}

	var total uint64 = 0
	for i, address := range toAddresses {
		if !wal.IsValidAddress(address) || address[:2] != "EC" {
			return trans, nil, fmt.Errorf("Invalid address given")
		}
		amt := rate * amounts[i]
		wal.Wallet.AddECOutput(trans, address, amt)
		total += amt
	}

	// Decide what addresses to pay with
	// Pay with largest first
	faAddresses, err := wal.Wallet.GetAllFCTAddresses()
	if err != nil {
		return trans, nil, err
	}

	var balances []AddressBalancePair

	for _, address := range faAddresses {
		addr := address.String()
		balance, err := factom.GetFactoidBalance(addr)
		if err != nil {
			return trans, nil, err
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
			return trans, nil, fmt.Errorf("Not enough factoids to cover the transaction")
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
		return trans, nil, err
	}

	// The last addresse we used to pay, we need to check if it can cover the fee
	if list[i].Balance < fee { // If it cannot, lets find one that can
		oldI := i
		i, err = checkForAddressForFee(list, transStruct, i, rate)
		if i == -1 || err != nil { // We don't have an address that can pay for the fee.
			return trans, nil, fmt.Errorf("Not enough factoids to cover the transaction")
		} else {
			if i != oldI { // It shouldn't, but doesn't hurt to double check
				err = wal.Wallet.AddInput(trans, list[i].Address, 0)
				if err != nil {
					return trans, nil, err
				}
			}
		}
	}

	err = wal.Wallet.AddFee(trans, list[i].Address, rate)
	if err != nil {
		return trans, nil, err
	}

	err = wal.Wallet.SignTransaction(trans)
	if err != nil {
		return trans, nil, err
	}

	r := new(ReturnTransStruct)
	r.Total = total
	r.Fee = fee

	return trans, r, nil
}

// A lot of parameters. This function is reused for EC and FCT transations. All it does it, if the last address input cannnot cover the fee
// this finds an address that can.
//	Parameters:
//		list = List of addresses
//		transStruct = The transaction structure that can calculate a fee
//		i = Last address in the list we have inputted into the transation
//		rate = current fee rate
func checkForAddressForFee(list []AddressBalancePair, transStruct *factoid.Transaction, i int, rate uint64) (indexToPay int, err error) {
	if i >= len(list) { // Out of addresses? Sorry, no transaction
		return -1, fmt.Errorf("Not enough factoids to cover the transaction")
	}

	// Ok, how much would another address input increase the fee to?
	fakeAddr := factoid.NewAddress(primitives.Sha([]byte("A fake address")).Bytes())
	transStruct.AddInput(fakeAddr, 0)          // We do not use this address, just need a new fee calc
	fee, err := transStruct.CalculateFee(rate) // We have the rate from earlier
	if err != nil {
		return -1, err
	}

	// So that is the fee, lets check if we have anything
	for list[i].Balance < fee {
		i++
		if i >= len(list) {
			return -1, fmt.Errorf("Not enough factoids to cover the transaction")
		}
	}

	// Sweet, we got one. Lets return the index of it
	return i, nil
}

// Initial attempt.
/*// The last addresse we used to pay, we need to check if it can cover the fee
if list[i].Balance < fee { // If it cannot, lets find one that can
	if i >= len(list) { // Out of addresses? Sorry, no transaction
		return trans, nil, fmt.Errorf("Not enough factoids to cover the transaction")
	}

	fakeAddr := factoid.NewAddress(primitives.Sha([]byte("A fake address")).Bytes())
	// Ok, how much would another address input increase the fee to?
	transStruct.AddInput(fakeAddr, 0)         // We may not use this address, just need a new fee calc
	fee, err = transStruct.CalculateFee(rate) // We have the rate from earlier
	if err != nil {
		return trans, nil, err
	}

	// So that is the fee, lets check if we have anything
	for list[i].Balance < fee {
		i++
		if i >= len(list) {
			return trans, nil, fmt.Errorf("Not enough factoids to cover the transaction")
		}
	}

	// Lets add it to the transaction and continue
	wal.Wallet.AddInput(trans, list[i].Address, 0)
}*/

func (wal *WalletDB) GetAddressBalance(address string) (uint64, error) {
	bal, err := factom.GetFactoidBalance(address)
	return uint64(bal), err
}

func (wal *WalletDB) SendTransaction(trans string) (string, error) {
	transObj, err := factom.SendTransaction(trans)
	if err != nil {
		return "", err
	}
	return transObj.TxID, nil
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
