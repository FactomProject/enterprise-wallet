package wallet_test

import (
	"testing"
	"time"

	. "github.com/FactomProject/enterprise-wallet/wallet"
)

func TestDisplayTransaction(t *testing.T) {
	a := newDisplayTransaction()
	b := newDisplayTransaction()

	if !a.IsSameAs(*b) {
		t.Fatal("Should be same, but are not")
	}

	// Make many different
	list := make([]*DisplayTransaction, 10)
	for i := range list {
		list[i] = newDisplayTransaction()
	}
	list[0].Inputs = append(list[0].Inputs, *newTransactionAddressInfo())
	list[1].Outputs = append(list[0].Outputs, *newTransactionAddressInfo())
	list[2].Inputs[0].Amount = 2
	list[3].Outputs[0].Amount = 2
	list[4].TotalInput = 2
	list[5].TotalFCTOutput = 2
	list[6].TotalECOutput = 2
	list[7].TxID = "different"
	list[8].Height = 100
	list[9].Action[0] = true

	for i := range list {
		if a.IsSameAs(*list[i]) {
			t.Fatal("Should not be same, but are")
		}
	}
}

func TestTransactionAddressInfo(t *testing.T) {
	a := NewTransactionAddressInfo("random", "add", 0, "fct")
	b := NewTransactionAddressInfo("random", "add", 0, "fct")

	if !a.IsSameAs(*b) {
		t.Fatal("Should be same, but are not")
	}

	b.Name = "notsame"
	if a.IsSameAs(*b) {
		t.Fatal("Should not be same, but are")
	}
}

func newDisplayTransaction() *DisplayTransaction {
	dt := new(DisplayTransaction)
	//dt.ITrans = t
	dt.TotalInput = 0
	dt.TotalFCTOutput = 0
	dt.TotalECOutput = 0
	dt.Height = 0
	dt.TxID = "Random"
	dt.Inputs = make([]TransactionAddressInfo, 1)
	dt.Inputs[0] = *newTransactionAddressInfo()
	dt.Outputs = make([]TransactionAddressInfo, 1)
	dt.Outputs[0] = *newTransactionAddressInfo()
	dt.TotalFCTOutput = 1
	dt.TotalECOutput = 0
	dt.TotalInput = 1
	dt.Action = [3]bool{false, false, false}
	dt.ExactTime = time.Now()
	dt.Date = dt.ExactTime.Format(("01/02/2006"))
	dt.Time = dt.ExactTime.Format(("15:04:05"))

	return dt
}

func newTransactionAddressInfo() *TransactionAddressInfo {
	a := NewTransactionAddressInfo("Name", "FA2Ax443J2xK63E38znfxrZ6kaVFmcfV7Hfy1Uj9dYV8LPB8m3Zf", 1, "FCT")
	return a
}
