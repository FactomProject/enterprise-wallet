package wallet

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// CheckTransactionAndGetName doublechecks the transaction is the same (with amounts and addresses)
// This is to confirm an already constructed transaction
func (wal *WalletDB) CheckTransactionAndGetName(toAddresses []string, amounts []string, feeAddress string) (string, error) {
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
		amts, err := StringAmountsToUin64Amounts(toAddresses, amounts)
		if err != nil {
			return "", err
		}
		for i, o := range outs {
			amt := amts[i]

			compAddr := ""
			if toAddresses[i][:2] == "FA" {
				compAddr = primitives.ConvertFctAddressToUserStr(o.GetAddress())
				if o.GetAmount() != uint64(amt) {
					if toAddresses[i] != feeAddress {
						return name, fmt.Errorf("A change in the amount of an output has been detected")
					}
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

// StringAmountsToUin64Amounts assumes amounts to be a float for a factoid and a uint64 for an entry credit
// Will multiply by 1e8 for factoids so "1" is 1 factoid. Not 1 factoshi. This is because this call usally
// comes from input from the user
func StringAmountsToUin64Amounts(addresses []string, amounts []string) ([]uint64, error) {
	var amts []uint64
	if len(addresses) != len(amounts) {
		return nil, fmt.Errorf("Length of addresses and amounts do not match")
	}
	for i, a := range amounts {
		if len(addresses[i]) < 20 {
			return nil, fmt.Errorf("Invalid address given")
		}
		if addresses[i][:2] == "FA" {
			amt64, err := strconv.ParseFloat(a, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid amount given")
			}
			amts = append(amts, uint64(amt64*1e8))
		} else {
			amt64, err := strconv.ParseUint(a, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid amount given")
			}
			amts = append(amts, uint64(amt64))
		}
	}

	return amts, nil
}

// CalculateNeededInput calculates how many factoids are needed to cover the outputs. Takes into consideration
// the EC rate if EC is output
func (wal *WalletDB) CalculateNeededInput(toAddresses []string, toAmounts []string) (uint64, error) {
	var toAmts []uint64
	toAmts, err := StringAmountsToUin64Amounts(toAddresses, toAmounts)
	if err != nil {
		return 0, err
	}

	rate, err := wal.GetRate()
	if err != nil {
		return 0, fmt.Errorf("Could not get the rate for converting entry credits. Factomd may be down or on a different port.\n")
	}

	var total uint64 = 0
	for i, address := range toAddresses {
		if !wal.IsValidAddress(address) {
			return 0, fmt.Errorf(" %s is not a valid address\n", address)
		}
		if address[:2] == "FA" {
			total += toAmts[i]
		} else if address[:2] == "EC" {
			total += toAmts[i] * rate
		} else {
			return 0, fmt.Errorf(" %s is not a valid public address\n", address)
		}
	}

	return total, nil
}

// ConstructTransactionFromValuesStrings constructs a transaction if all inputs are already given
// Amounts are parsed into a float or uint64 depending on factoid/ec
func (wal *WalletDB) ConstructTransactionFromValuesStrings(toAddresses []string, toAmounts []string, fromAddresses []string, fromAmounts []string, feeAddress string, sign bool) (string, *ReturnTransStruct, error) {
	if len(toAddresses) != len(toAmounts) {
		return "", nil, fmt.Errorf("Lengths of output addresses to amounts does not match")
	} else if len(fromAddresses) != len(fromAmounts) {
		return "", nil, fmt.Errorf("Lengths of input addresses to amounts does not match")
	}
	var err error

	var toAmts []uint64
	toAmts, err = StringAmountsToUin64Amounts(toAddresses, toAmounts)
	if err != nil {
		return "", nil, err
	}

	var fromAmts []uint64
	fromAmts, err = StringAmountsToUin64Amounts(fromAddresses, fromAmounts)
	if err != nil {
		return "", nil, err
	}

	return wal.ConstructTransactionFromValues(toAddresses, toAmts, fromAddresses, fromAmts, feeAddress, sign)
}

// ConstructTransactionFromValues constructs a transaction from given input and output values. An error might contain the amount of input needed aswell if it is incorrect
func (wal *WalletDB) ConstructTransactionFromValues(toAddresses []string, toAmounts []uint64, fromAddresses []string, fromAmounts []uint64, feeAddress string, sign bool) (string, *ReturnTransStruct, error) {
	if len(toAddresses) != len(toAmounts) {
		return "", nil, fmt.Errorf("Lengths of output addresses to amounts does not match")
	} else if len(fromAddresses) != len(fromAmounts) {
		return "", nil, fmt.Errorf("Lengths of input addresses to amounts does not match")
	} else if !(wal.IsValidAddress(feeAddress) && feeAddress[:2] == "FA") {
		return "", nil, fmt.Errorf("Invalid address for fee")
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

	rate, err := wal.GetRate()
	if err != nil {
		return trans, nil, err
	}

	var total uint64 = 0
	var amt uint64
	for i, address := range toAddresses {
		if !wal.IsValidAddress(address) {
			return trans, nil, fmt.Errorf("Invalid address given")
		}
		if toAddresses[i][:2] == "FA" {
			amt = toAmounts[i]
			err = wal.Wallet.AddOutput(trans, address, amt)
		} else if toAddresses[i][:2] == "EC" {
			amt = rate * toAmounts[i]
			err = wal.Wallet.AddECOutput(trans, address, amt)
		} else {
			return trans, nil, fmt.Errorf(address + " is not a public address")
		}
		if err != nil {
			return trans, nil, err
		}
		total += amt
	}

	var totalIn uint64 = 0
	for _, a := range fromAmounts {
		totalIn += a
	}

	var feeAddIndex int = -1
	var feeAddBal uint64 = 0
	for i, address := range fromAddresses {
		addBal, err := wal.GetAddressBalance(address)
		if err != nil {
			return trans, nil, err
		}

		if fromAmounts[i] > addBal {
			return trans, nil, fmt.Errorf("%s only has %s FCT, and cannot cover the %s FCT input."+
				" If you are sure this balance is incorrect, make sure factomd is synced.",
				address,
				strconv.FormatFloat(float64(addBal)/1e8, 'f', -1, 64),
				strconv.FormatFloat(float64(fromAmounts[i])/1e8, 'f', -1, 64))
		}

		// for later use
		if fromAddresses[i] == feeAddress {
			feeAddIndex = i
			feeAddBal = addBal
		}

		err = wal.Wallet.AddInput(trans, address, fromAmounts[i])
		if err != nil {
			return trans, nil, err
		}
	}

	if total > totalIn {
		return trans, nil, fmt.Errorf("The amount of input is not enough to cover the transaction. The needed input is: %f FCT.\n", float64(total)/1e8)
	} else if total < totalIn {
		return trans, nil, fmt.Errorf("The amount of input is too much for the transaction. The needed input is: %f FCT.\n", float64(total)/1e8)
	}

	transStruct := wal.Wallet.GetTransactions()[trans]
	if transStruct == nil {
		return trans, nil, fmt.Errorf("Transaction not found")
	}

	fee, err := transStruct.CalculateFee(rate)
	if err != nil {
		return trans, nil, err
	}

	feeTakenCareOf := false // Did the loop find an address to deduct from
	for _, add := range toAddresses {
		if add[:2] == "FA" {
			if add == feeAddress {
				err := wal.Wallet.SubFee(trans, add, rate)
				if err != nil {
					return trans, nil, err
				}
				feeTakenCareOf = true
				break
			}
		}
	}

	if !feeTakenCareOf {
		if feeAddIndex == -1 {
			return trans, nil, fmt.Errorf("An error occured while adding the fee")
		}

		if fee+fromAmounts[feeAddIndex] > feeAddBal {
			return trans, nil, fmt.Errorf("The inputs can cover the outputs, but the fee address selected "+
				"cannot cover the fee. The total input from %s is %s FCT, but it only has %s FCT. Choose another "+
				"input to cover the fee, or choose an output to cover the fee.",
				fromAddresses[feeAddIndex],
				strconv.FormatFloat(float64(feeAddBal)/1e8, 'f', -1, 64),
				strconv.FormatFloat(float64(fromAmounts[feeAddIndex])/1e8, 'f', -1, 64))
		}

		err = wal.Wallet.AddFee(trans, feeAddress, rate)
		if err != nil {
			return trans, nil, err
		}
	}

	if sign {
		err = wal.Wallet.SignTransaction(trans, true)
		if err != nil {
			return trans, nil, err
		}
	}

	r := new(ReturnTransStruct)
	r.Total = total
	r.Fee = fee

	return trans, r, nil
}

func (wal *WalletDB) ConstructSendFactoidsStrings(toAddresses []string, amounts []string) (string, *ReturnTransStruct, error) {
	var amts []uint64
	amts, err := StringAmountsToUin64Amounts(toAddresses, amounts)
	if err != nil {
		return "", nil, err
	}

	return wal.ConstructTransaction(toAddresses, amts)
}

func (wal *WalletDB) ConstructConvertEntryCreditsStrings(toAddresses []string, amounts []string) (string, *ReturnTransStruct, error) {
	var amts []uint64
	amts, err := StringAmountsToUin64Amounts(toAddresses, amounts)
	if err != nil {
		return "", nil, err
	}

	return wal.ConstructTransaction(toAddresses, amts)
}

func (wal *WalletDB) ImportTransaction(name string, hex string) error {
	return wal.Wallet.ImportComposedTransaction(name, hex)
}

func (wal *WalletDB) ExportTransaction(name string) (string, error) {
	req, err := wal.Wallet.ComposeTransaction(name)
	if err != nil {
		return "", err
	}

	return req.String(), nil
}

func (wal *WalletDB) DeleteTransaction(trans string) error {
	return wal.Wallet.DeleteTransaction(trans)
}

// ConstructTransaction
// Transaction name is hash of all the addresses. More than 1 transaction to
// an address(es) should not be open, but combined.
// The output is determined by the output address for ECOutput or FCTOutput
// Parameters:
//		toAddresses = list of output addresses
//		amounts = list of amounts to each output, indicies must match
// Returns:
//		Transaction Name, Transaction Info, error

func (wal *WalletDB) ConstructTransaction(toAddresses []string, amounts []uint64) (string, *ReturnTransStruct, error) {
	if len(toAddresses) != len(amounts) {
		return "", nil, fmt.Errorf("Lengths of address to amount does not match")
	} else if len(toAddresses) == 0 {
		return "", nil, fmt.Errorf("No recipient given")
	}

	trans := hashStringList(toAddresses)
	trans = trans[:32] // Name of transaction

	// If the transaction exists, we will overwrite it
	transMap := wal.Wallet.GetTransactions()
	if t, _ := transMap[trans]; t != nil {
		wal.DeleteTransaction(trans)
	}

	err := wal.Wallet.NewTransaction(trans)
	if err != nil {
		return trans, nil, err
	}

	rate, err := wal.GetRate()
	if err != nil {
		return trans, nil, err
	}

	var total uint64 = 0
	var amt uint64
	for i, address := range toAddresses {
		if !wal.IsValidAddress(address) {
			return trans, nil, fmt.Errorf("Invalid address given")
		}
		if toAddresses[i][:2] == "FA" {
			amt = amounts[i]
			err = wal.Wallet.AddOutput(trans, address, amt)
		} else if toAddresses[i][:2] == "EC" {
			amt = rate * amounts[i]
			err = wal.Wallet.AddECOutput(trans, address, amt)
		} else {
			return trans, nil, fmt.Errorf(address + " is not a public address")
		}
		if err != nil {
			return trans, nil, err
		}
		total += amt
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
		balance, err := wal.GetAddressBalance(addr)
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
			err := wal.Wallet.AddInput(trans, list[i].Address, totalLeft)
			if err != nil {
				return trans, nil, err
			}
			list[i].Balance -= totalLeft
			totalLeft = 0
		} else {
			if list[i].Balance > 0 {
				err := wal.Wallet.AddInput(trans, list[i].Address, list[i].Balance)
				if err != nil {
					return trans, nil, err
				}
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
			err := wal.Wallet.AddInput(trans, list[i].Address, 0)
			if err != nil {
				return trans, nil, err
			}
		}
	}

	err = wal.Wallet.AddFee(trans, list[i].Address, rate)
	if err != nil {
		return trans, nil, err
	}

	err = wal.Wallet.SignTransaction(trans, true)
	if err != nil {
		return trans, nil, err
	}

	r := new(ReturnTransStruct)
	r.Total = total
	r.Fee = fee

	return trans, r, nil
}

// checkForAddressForFee: A lot of parameters. This function is reused for EC and FCT transactions. All it does it,
// if the last address input cannot cover the fee this finds an address that can.
//	Parameters:
//		list = List of addresses
//		transStruct = The transaction structure that can calculate a fee
//		i = Last address in the list we have inputted into the transaction
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

// GetAddressBalance Return balance for an address
// Checks the balance cache first, if nothing is found it
// sends a request to factomd
func (wal *WalletDB) GetAddressBalance(address string) (uint64, error) {
	if bal, hit := wal.BalanceCache.Get(address); hit {
		return bal, nil
	}
	bal, err := factom.GetFactoidBalance(address)
	if err == nil {
		wal.BalanceCache.Set(address, uint64(bal))
	}
	return uint64(bal), err
}

// GetRate Return current EC rate,
// Cached for 60 seconds
func (wal *WalletDB) GetRate() (uint64, error) {
	if rate, hit := wal.BalanceCache.Get("rate"); hit {
		return rate, nil
	}

	rate, err := factom.GetRate()
	if err != nil {
		return 0, err
	}

	wal.BalanceCache.Set("rate", rate)
	return rate, nil
}

type SendTransactionResp struct {
	Message string `json:"message"`
	Txid    string `json:"txid"`
}

func (wal *WalletDB) SendTransaction(trans string) (string, error) {
	req, err := wal.Wallet.ComposeTransaction(trans)
	if err != nil {
		return "", err
	}

	respJson, err := factom.SendFactomdRequest(req)
	if err != nil {
		return "", err
	}

	resp := new(SendTransactionResp)
	result, err := respJson.Result.MarshalJSON()
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(result, resp)
	if err != nil {
		return "", err
	}
	return resp.Txid, nil
}

func (wal *WalletDB) FactoidAddressToHumanReadable(add interfaces.IAddress) string {
	return primitives.ConvertFctAddressToUserStr(add)
}

func FactoidAddressToHumanReadable(add interfaces.IAddress) string {
	return primitives.ConvertFctAddressToUserStr(add)
}

func (wal *WalletDB) ECAddressToHumanReadable(add interfaces.IAddress) string {
	return primitives.ConvertECAddressToUserStr(add)
}

func ECAddressToHumanReadable(add interfaces.IAddress) string {
	return primitives.ConvertECAddressToUserStr(add)
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
