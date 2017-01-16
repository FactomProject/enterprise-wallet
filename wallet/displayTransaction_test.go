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

	b.Action[0] = true
	if a.IsSameAs(*b) {
		t.Fatal("Should not be same, but are")
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
	dt.Inputs = make([]TransactionAddressInfo, 0)
	dt.Outputs = make([]TransactionAddressInfo, 0)
	dt.Action = [3]bool{false, false, false}
	dt.ExactTime = time.Now()
	dt.Date = dt.ExactTime.Format(("01/02/2006"))
	dt.Time = dt.ExactTime.Format(("15:04:05"))

	return dt
}
