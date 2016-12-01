package address

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/FactomProject/btcutil/base58"
	"github.com/FactomProject/factom"
)

var _ = fmt.Sprintf("")

const maxNameLength int = 20

// Name/Address pair
type AddressNamePair struct {
	Name    string // Length maxNameLength Characters
	Address string // Length 52 Characters

	// Not Marshaled
	Balance float64 // Unused except for JSON return
}

func NewAddress(name string, address string) (*AddressNamePair, error) {
	if len(name) > maxNameLength {
		return nil, fmt.Errorf("Name must be max %d characters", maxNameLength)
	} else if len(address) != 52 {
		return nil, errors.New("Address must be 52 characters")
	}

	if !factom.IsValidAddress(address) {
		return nil, errors.New("Address is invalid")
	}

	add := new(AddressNamePair)

	//var n [maxNameLength]byte
	//copy(n[:maxNameLength], name)

	add.Name = name
	add.Address = address

	return add, nil
}

func (anp *AddressNamePair) ChangeName(name string) error {
	if len(name) > maxNameLength {
		return fmt.Errorf("Name too long, must be less than %d characters", maxNameLength)
	}
	anp.Name = name
	return nil
}

func (anp *AddressNamePair) IsSameAs(b *AddressNamePair) bool {
	if strings.Compare(anp.Name, b.Name) == 0 {
		if strings.Compare(anp.Address, b.Address) == 0 {
			return true
		}
	}

	return false
}

func (anp *AddressNamePair) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)

	var n [maxNameLength]byte
	copy(n[:maxNameLength], anp.Name)
	buf.Write(n[:maxNameLength])

	add := base58.Decode(anp.Address)
	var a [38]byte
	copy(a[:38], add[:])

	buf.Write(a[:38])

	return buf.Next(buf.Len()), nil
}

func (anp *AddressNamePair) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	newData = data

	nameData := bytes.Trim(newData[:maxNameLength], "\x00")
	anp.Name = fmt.Sprintf("%s", nameData)
	newData = newData[maxNameLength:]

	anp.Address = base58.Encode(newData[:38])

	if len(newData) > 38 {
		newData = newData[38:]
	} else {
		newData = nil
	}

	return
}

func (anp *AddressNamePair) UnmarshalBinary(data []byte) (err error) {
	_, err = anp.UnmarshalBinaryData(data)
	return
}

//Address List
type AddressList struct {
	Length uint32
	List   []AddressNamePair
}

func NewAddressList() *AddressList {
	addList := new(AddressList)
	addList.Length = 0

	return addList
}

// Searches for Address
func (addList *AddressList) Get(address string) (*AddressNamePair, int) {
	if len(address) != 52 {
		return nil, -1
	}

	for i, ianp := range addList.List {
		if strings.Compare(ianp.Address, address) == 0 {
			return &ianp, i
		}
	}
	return nil, -1
}

func (addList *AddressList) AddANP(anp *AddressNamePair) error {
	if len(anp.Name) == 0 || len(anp.Address) != 52 {
		return errors.New("Nil AddressNamePair")
	}

	_, i := addList.Get(anp.Address)
	if i == -1 {
		addList.List = append(addList.List, *anp)
		addList.Length++
		return nil
	}

	// Duplicate Found
	return errors.New("Address or Name already exists")

}

func (addList *AddressList) Add(name string, address string) (*AddressNamePair, error) {
	// We check for valid factom address higher up, this is just a basic check
	if len(name) == 0 || len(address) != 52 {
		return nil, errors.New("Nil AddressNamePair")
	}

	anp, err := NewAddress(name, address)
	if err != nil {
		return nil, err
	}

	_, i := addList.Get(anp.Address)
	if i == -1 {
		addList.List = append(addList.List, *anp)
		addList.Length++
		return anp, nil
	}

	// Duplicate Found
	return nil, errors.New("Address already exists")

}

func (addList *AddressList) Remove(removeAnp *AddressNamePair) error {
	_, i := addList.Get(removeAnp.Address)
	if i == -1 {
		return errors.New("Not found")
	}
	addList.Length--

	addList.List = append(addList.List[:i], addList.List[i+1:]...)
	return nil
}

func (addList *AddressList) RemoveAddress(removeAdd string) error {
	ranp, err := NewAddress("", removeAdd)
	if err != nil {
		return err
	}

	return addList.Remove(ranp)
}

func (addList *AddressList) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	var number [8]byte
	binary.BigEndian.PutUint32(number[:], addList.Length)
	buf.Write(number[:])

	for _, anp := range addList.List {
		anpData, err := anp.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(anpData)
	}

	return buf.Next(buf.Len()), err
}

func (addList *AddressList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	newData = data

	addList.Length = binary.BigEndian.Uint32(data[:8])
	newData = newData[8:]

	var i uint32 = 0
	for i < addList.Length {
		anp := new(AddressNamePair)
		newData, err = anp.UnmarshalBinaryData(newData)
		addList.List = append(addList.List, *anp)
		i++
	}

	return
}

func (addList *AddressList) UnmarshalBinary(data []byte) (err error) {
	_, err = addList.UnmarshalBinaryData(data)
	return
}

func (addList *AddressList) IsSameAs(b *AddressList) bool {
	if addList.Length != b.Length {
		return false
	}

	for _, anp := range addList.List {
		if inap, i := b.Get(anp.Address); i == -1 || !anp.IsSameAs(inap) {
			return false
		}
	}

	return true
}
