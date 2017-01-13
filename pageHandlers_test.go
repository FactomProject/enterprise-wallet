package main_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	. "github.com/FactomProject/enterprise-wallet"
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
