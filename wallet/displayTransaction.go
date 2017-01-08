package wallet

import (
	//"bytes"
	//"encoding/binary"
	//"fmt"
	"time"
	//"github.com/FactomProject/enterprise-wallet/address"
	//"github.com/FactomProject/btcutil/base58"
)

// Names are "" if not in wallet
type DisplayTransaction struct {
	Inputs     []TransactionAddressInfo
	TotalInput uint64

	Outputs        []TransactionAddressInfo
	TotalFCTOutput uint64
	TotalECOutput  uint64

	TxID      string
	Height    uint32
	Action    [3]bool // Sent, recieved, converted
	Date      string
	Time      string
	ExactTime time.Time

	//ITrans interfaces.ITransaction
}

func (a *DisplayTransaction) IsSameAs(b DisplayTransaction) bool {
	if !a.IsSimilarTo(b) {
		return false
	}

	for i := 0; i < 3; i++ {
		if a.Action[i] != b.Action[i] {
			return false
		}
	}

	return true
}

// Does not count actions
func (a *DisplayTransaction) IsSimilarTo(b DisplayTransaction) bool {
	if len(a.Inputs) != len(b.Inputs) {
		return false
	}
	if len(a.Outputs) != len(b.Outputs) {
		return false
	}

	for i := 0; i < len(a.Inputs); i++ {
		if !a.Inputs[i].IsSimilarTo(b.Inputs[i]) {
			return false
		}
	}
	if a.TotalInput != b.TotalInput {
		return false
	}
	for i := 0; i < len(a.Outputs); i++ {
		if !a.Outputs[i].IsSimilarTo(b.Outputs[i]) {
			return false
		}
	}
	if a.TotalFCTOutput != b.TotalFCTOutput {
		return false
	}
	if a.TotalECOutput != b.TotalECOutput {
		return false
	}
	if a.TxID != b.TxID {
		return false
	}
	if a.Height != b.Height {
		return false
	}

	return true
}

/* TransactionAddressInfo */
type TransactionAddressInfo struct {
	Name    string
	Address string
	Amount  uint64
	Type    string // FCT or EC
}

/* Not Needed
func (a *TransactionAddressInfo) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	var n [address.MaxNameLength]byte
	copy(n[:address.MaxNameLength], a.Name)
	buf.Write(n[:address.MaxNameLength])

	add := base58.Decode(a.Address)
	var b [38]byte
	copy(b[:38], add[:])
	buf.Write(b[:38])

	var number [8]byte
	binary.BigEndian.PutUint64(number[:], a.Amount)
	buf.Write(number[:])

	var t [3]byte
	copy(t[:3], a.Type)
	buf.Write(t[:3])

	return buf.Next(buf.Len()), nil
} */

func NewTransactionAddressInfo(name string, address string, amount uint64, tokenType string) *TransactionAddressInfo {
	t := new(TransactionAddressInfo)
	t.Name = name
	t.Address = address
	t.Amount = amount
	t.Type = tokenType

	return t
}

func (a *TransactionAddressInfo) IsSimilarTo(b TransactionAddressInfo) bool {
	if a.Address != b.Address {
		return false
	}
	if a.Amount != b.Amount {
		return false
	}
	if a.Type != b.Type {
		return false
	}
	return true
}

func (a *TransactionAddressInfo) IsSameAs(b TransactionAddressInfo) bool {
	if !a.IsSimilarTo(b) {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	return true
}
