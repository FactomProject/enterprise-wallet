package wallet

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/FactomProject/M2GUIWallet/address"
	"github.com/FactomProject/factom"
)

// Wallet use outside DB
type WalletStruct struct {
	FactoidAddresses     *address.AddressList
	EntryCreditAddresses *address.AddressList
	ExternalAddresses    *address.AddressList

	// Not marshaled into database
	FactoidTotal int64
	ECTotal      int64

	sync.RWMutex
}

func NewWallet() *WalletStruct {
	w := new(WalletStruct)
	w.FactoidAddresses = address.NewAddressList()
	w.EntryCreditAddresses = address.NewAddressList()
	w.ExternalAddresses = address.NewAddressList()

	return w
}

func (w *WalletStruct) AddAddress(name string, address string, list int) (*address.AddressNamePair, error) {
	if err := w.addAddress(name, address, list); err != nil {
		return nil, err
	}

	w.Lock()
	defer w.Unlock()

	switch list {
	case 1:
		return w.FactoidAddresses.Add(name, address)
	case 2:
		return w.EntryCreditAddresses.Add(name, address)
	case 3:
		return w.ExternalAddresses.Add(name, address)
	}

	return nil, fmt.Errorf("Encountered an error, this should not be able to happen")
}

func (w *WalletStruct) AddSeededAddress(name string, address string, list int) (*address.AddressNamePair, error) {
	if err := w.addAddress(name, address, list); err != nil {
		return nil, err
	}

	w.Lock()
	defer w.Unlock()

	switch list {
	case 1:
		return w.FactoidAddresses.AddSeeded(name, address)
	case 2:
		return w.EntryCreditAddresses.AddSeeded(name, address)
	case 3:
		return w.ExternalAddresses.Add(name, address) // You can't do this, but will not hurt
	}

	return nil, fmt.Errorf("Encountered an error, this should not be able to happen")
}

func (w *WalletStruct) addAddress(name string, address string, list int) error {
	if list > 3 || list <= 0 {
		return fmt.Errorf("Invalid list")
	}

	if list == 3 {
		anp, getList, _ := w.GetAddress(address)
		if getList != -1 {
			return fmt.Errorf("You cannot add this address as it is located in your Addressbook. " +
				"It's nickname is: " + anp.Name)
		}
	}

	switch list {
	case 1: // Factoid
		if address[:2] != "FA" {
			return fmt.Errorf("Not a valid factoid address")
		}
	case 2: // EC
		if address[:2] != "EC" {
			return fmt.Errorf("Not a valid entry credit address")
		}
	case 3: // Either
		if !(address[:2] == "EC" || address[:2] == "FA") {
			return fmt.Errorf("Not a valid address")
		}
	}

	valid := factom.IsValidAddress(address)
	if !valid {
		return fmt.Errorf("Not a valid address")
	}

	return nil
}

func (w *WalletStruct) GetTotalAddressCount() uint64 {
	w.RLock()
	defer w.RUnlock()
	return w.FactoidAddresses.Length + w.EntryCreditAddresses.Length + w.ExternalAddresses.Length
}

// List is -1 for not found, 1 for FactoidAddressList, 2 for EntryCreditList, 3 for External
func (w *WalletStruct) GetAddress(address string) (anp *address.AddressNamePair, list int, index int) {
	w.RLock()
	defer w.RUnlock()

	list = -1

	anp, index = w.FactoidAddresses.Get(address)
	if index != -1 && anp != nil {
		list = 1
		return
	}

	anp, index = w.EntryCreditAddresses.Get(address)
	if index != -1 && anp != nil {
		list = 2
		return
	}

	anp, index = w.ExternalAddresses.Get(address)
	if index != -1 && anp != nil {
		list = 3
		return
	}

	return
}

func (w *WalletStruct) ChangeAddressName(address string, toName string) error {
	anp, list, i := w.GetAddress(address)
	if list == -1 || anp == nil || i == -1 {
		return fmt.Errorf("Address not found")
	}

	w.Lock()
	defer w.Unlock()
	if strings.Compare(anp.Address, address) == 0 { // To be sure
		switch list {
		case 1:
			//w.FactoidAddresses.List[i].Name = toName
			err := w.FactoidAddresses.List[i].ChangeName(toName)
			if err != nil {
				return err
			}
			return nil
		case 2:
			err := w.EntryCreditAddresses.List[i].ChangeName(toName)
			if err != nil {
				return err
			}
			return nil
		case 3:
			err := w.ExternalAddresses.List[i].ChangeName(toName)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("Could not change name")
}

func (w *WalletStruct) GetAllAddresses() []address.AddressNamePair {
	w.RLock()
	defer w.RUnlock()
	var anpList []address.AddressNamePair
	anpList = append(anpList, w.FactoidAddresses.List...)
	anpList = append(anpList, w.EntryCreditAddresses.List...)
	anpList = append(anpList, w.ExternalAddresses.List...)

	return anpList
}

func (w *WalletStruct) GetAllMyGUIAddresses() []address.AddressNamePair {
	w.RLock()
	defer w.RUnlock()
	var anpList []address.AddressNamePair
	anpList = append(anpList, w.FactoidAddresses.List...)
	anpList = append(anpList, w.EntryCreditAddresses.List...)

	return anpList
}

func (w *WalletStruct) GetAllAddressesFromList(list int) []address.AddressNamePair {
	w.RLock()
	defer w.RUnlock()
	var anpList []address.AddressNamePair
	switch list {
	case 1:
		anpList = append(anpList, w.FactoidAddresses.List...)
	case 2:
		anpList = append(anpList, w.EntryCreditAddresses.List...)
	case 3:
		anpList = append(anpList, w.ExternalAddresses.List...)
	}

	return anpList
}

func (w *WalletStruct) IsSameAs(b *WalletStruct) bool {
	w.RLock()
	defer w.RUnlock()
	b.RLock()
	defer b.RUnlock()

	if !w.FactoidAddresses.IsSameAs(b.FactoidAddresses) {
		return false
	} else if !w.EntryCreditAddresses.IsSameAs(b.EntryCreditAddresses) {
		return false
	} else if !w.ExternalAddresses.IsSameAs(b.ExternalAddresses) {
		return false
	}
	return true
}

func (w *WalletStruct) MarshalBinary() ([]byte, error) {
	w.RLock()
	defer w.RUnlock()
	buf := new(bytes.Buffer)

	data, err := w.FactoidAddresses.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = w.EntryCreditAddresses.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = w.ExternalAddresses.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (w *WalletStruct) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	w.Lock()
	defer w.Unlock()

	if w.FactoidAddresses == nil {
		w.FactoidAddresses = address.NewAddressList()
	}
	if w.EntryCreditAddresses == nil {
		w.EntryCreditAddresses = address.NewAddressList()
	}
	if w.ExternalAddresses == nil {
		w.ExternalAddresses = address.NewAddressList()
	}

	newData = data
	newData, err = w.FactoidAddresses.UnmarshalBinaryData(newData)
	if err != nil {
		return
	}

	newData, err = w.EntryCreditAddresses.UnmarshalBinaryData(newData)
	if err != nil {
		return
	}

	newData, err = w.ExternalAddresses.UnmarshalBinaryData(newData)
	if err != nil {
		return
	}

	return
}

func (w *WalletStruct) UnmarshalBinary(data []byte) error {
	_, err := w.UnmarshalBinaryData(data)
	return err
}

func (w *WalletStruct) RemoveAddressFromAnyList(address string) (*address.AddressNamePair, error) {
	anp, list, _ := w.GetAddress(address)
	if list > 3 {
		return nil, fmt.Errorf("This should never happen")
	}
	_, err := w.RemoveAddress(address, list)
	if err != nil {
		return nil, err
	}
	return anp, nil
}

func (w *WalletStruct) RemoveAddress(address string, list int) (string, error) {
	w.Lock()
	defer w.Unlock()

	switch list {
	case 0:
		return "", fmt.Errorf("No address found")
	case 1:
		err := w.FactoidAddresses.Remove(address)
		if err != nil {
			return "", err
		}

		// factom-wallet remove?
		return address, nil
	case 2:
		err := w.EntryCreditAddresses.Remove(address)
		if err != nil {
			return "", err
		}

		// factom-wallet remove?
		return address, nil
	case 3:
		err := w.ExternalAddresses.Remove(address)
		if err != nil {
			return "", err
		}

		// factom-wallet remove?
		return address, nil
	}

	return "", fmt.Errorf("Impossible to reach.")
}

// Adds balances to addresses so the GUI can display
func (w *WalletStruct) AddBalancesToAddresses() {
	w.Lock()
	defer w.Unlock()

	w.FactoidTotal = 0
	w.ECTotal = 0

	if w.FactoidAddresses.Length > 0 {
		for i, fa := range w.FactoidAddresses.List {
			bal, err := factom.GetFactoidBalance(fa.Address)
			if err != nil {
				w.FactoidAddresses.List[i].Balance = -1
			} else {
				w.FactoidAddresses.List[i].Balance = float64(bal) / 1e8
				w.FactoidTotal += bal
			}
		}

		for i, ec := range w.EntryCreditAddresses.List {
			bal, err := factom.GetECBalance(ec.Address)
			if err != nil {
				w.EntryCreditAddresses.List[i].Balance = -1
			} else {
				w.EntryCreditAddresses.List[i].Balance = float64(bal)
				w.ECTotal += bal
			}
		}

		for i, a := range w.ExternalAddresses.List {
			if a.Address[:2] == "FA" {
				bal, err := factom.GetFactoidBalance(a.Address)
				if err != nil {
					w.ExternalAddresses.List[i].Balance = -1
				} else {
					w.ExternalAddresses.List[i].Balance = float64(bal) / 1e8
				}
			} else if a.Address[:2] == "EC" {
				bal, err := factom.GetECBalance(a.Address)
				if err != nil {
					w.ExternalAddresses.List[i].Balance = -1
				} else {
					w.ExternalAddresses.List[i].Balance = float64(bal)
				}
			}
		}
	}
}
