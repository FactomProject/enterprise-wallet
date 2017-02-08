package main_test

import (
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"testing"

	. "github.com/FactomProject/enterprise-wallet"
	"github.com/FactomProject/factomd/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestOldToNewUnmarshal(t *testing.T) {
	s := new(SettingsStruct)

	data, err := hex.DecodeString("66616c736566616c736566616c736566616c7365")
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = s.UnmarshalBinaryData(data)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestSettings(t *testing.T) {
	s := new(SettingsStruct)

	n, err := MarshalSettingAndGetNewUnmarshaled(s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !n.IsSameAs(s) {
		t.Fatal("Not the Same")
	}

	s.KeyExport = true
	n, err = MarshalSettingAndGetNewUnmarshaled(s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !n.IsSameAs(s) {
		t.Fatal("Not the Same")
	}

	s.ImportExport = true
	n, err = MarshalSettingAndGetNewUnmarshaled(s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !n.IsSameAs(s) {
		t.Fatal("Not the Same")
	}

	s.CoinControl = true
	n, err = MarshalSettingAndGetNewUnmarshaled(s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !n.IsSameAs(s) {
		t.Fatal("Not the Same")
	}

	s.FactomdLocation = "Another type of string"
	n, err = MarshalSettingAndGetNewUnmarshaled(s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !n.IsSameAs(s) {
		t.Fatal("Not the Same")
	}

	s.SetFactomdLocation("random")
}

func MarshalSettingAndGetNewUnmarshaled(a *SettingsStruct) (*SettingsStruct, error) {
	data, err := a.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("Did not marshal")
	}
	_ = data

	n := new(SettingsStruct)
	newdata, err := n.UnmarshalBinaryData(data)
	if err != nil {
		return nil, fmt.Errorf("Did not unmarshal", err)
	}
	if len(newdata) != 0 {
		return nil, fmt.Errorf("Did not unmarshal correctly")
	}

	return n, nil
}

func TestVersionChange(t *testing.T) {
	dataStr := []string{
		"66616c736566616c736566616c736566616c7365000000000000000000000000000000000000000000000000000000000000",
		"66616c7365747275650066616c736566616c7365000000000000000000000000000000000000000000000000000000000000",
		"66616c7365747275650066616c73657472756500000000000000000000000000000000000000000000000000000000000000",
		"66616c7365747275650074727565007472756500000000000000000000000000000000000000000000000000000000000000",
		"66616c7365747275650074727565007472756500416e6f746865722074797065206f6620737472696e670000000000000000",
		"66616c7365747275650074727565007472756500",
		"66616c736566616c736566616c736566616c736548656c6c6f2100"}

	for _, s := range dataStr {
		se := new(SettingsStruct)
		data, err := hex.DecodeString(s)
		if err != nil {
			t.Error(err)
		}

		data, err = se.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}

		if len(data) != 0 {
			t.Errorf("Should be length 0, found %d", len(data))
		}
	}

	for i := 0; i < 1000; i++ {
		se := new(SettingsStruct)
		str := random.RandomString()
		max := random.RandIntBetween(0, MAX_FACTOMDLOCATION_SIZE)
		if len(str) > max {
			str = str[:max]
		}

		se.FactomdLocation = str

		data, err := se.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		sa := new(SettingsStruct)
		data, err = sa.UnmarshalBinaryData(data)
		if !sa.IsSameAs(se) {
			t.Error("Should be same")
		}

		if len(data) != 0 {
			t.Errorf("Data length should be 0, found %d", len(data))
		}
	}

}

// Cannot really test to verify the data, will just test if they don't fail
func TestHandlers(t *testing.T) {
	MasterSettings = new(SettingsStruct)
	InitTemplate()
	r := httptest.NewRequest("GET", "localhost:8091", nil)
	w := httptest.NewRecorder()

	var err error
	err = HandleIndexPage(w, r)
	if err != nil {
		t.Fail()
	}

	err = HandleAddressBook(w, r)
	if err != nil {
		t.Fail()
	}

	err = HandleSettings(w, r)
	if err != nil {
		t.Fail()
	}

	/* Have to add form values
	err = HandleEditAddressFactoids(w, r)
	if err != nil {
		t.Error("Failed on Index Page:", err.Error())
	}*/

	err = HandleImportExportTransaction(w, r)
	if err != nil {
		t.Fail()
	}

	err = HandleNewAddress(w, r)
	if err != nil {
		t.Fail()
	}

	err = HandleNewAddressFactoid(w, r)
	if err != nil {
		t.Fail()
	}

	err = HandleNewAddressEntryCredits(w, r)
	if err != nil {
		t.Fail()
	}

	err = HandleNewAddressExternal(w, r)
	if err != nil {
		t.Fail()
	}

	/* Have to add form values
	err = HandleReceiveFactoids(w, r)
	if err != nil {
		t.Fail()
	}*/

	err = HandleSendFactoids(w, r)
	if err != nil {
		t.Fail()
	}

	err = HandleCreateEntryCredits(w, r)
	if err != nil {
		t.Fail()
	}

	err = HandleNotFoundError(w, r)
	if err != nil {
		t.Fail()
	}
}
