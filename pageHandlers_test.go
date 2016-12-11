package main_test

import (
	"fmt"
	"testing"

	. "github.com/FactomProject/M2GUIWallet"
)

var _ = fmt.Sprintf("")

func TestSettings(t *testing.T) {
	s := new(SettingsStruct)
	s.KeyExport = true

	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatalf("Did not marshal")
	}
	_ = data

	n := new(SettingsStruct)
	newdata, err := n.UnmarshalBinaryData(data)
	if err != nil {
		t.Fatalf("Did not unmarshal", err)
	}
	if len(newdata) != 0 {
		t.Fatalf("Did not unmarshal correctly")
	}

	if n.DarkTheme != s.DarkTheme {
		t.Fatalf("Does not match")
	} else if n.Theme != s.Theme {
		t.Fatalf("Does not match")
	} else if n.KeyExport != s.KeyExport {
		t.Fatalf("Does not match")
	} else if n.CoinControl != s.CoinControl {
		t.Fatalf("Does not match")
	}
}
