package main_test

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/FactomProject/enterprise-wallet"
	"github.com/FactomProject/factomd/common/primitives/random"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMarshal(t *testing.T) {
	for i := 0; i < 100; i++ {
		str := random.RandomString()
		max := random.RandIntBetween(0, 100)
		if max < len(str) {
			str = str[:max]
		}

		data, err := MarshalStringToBytes(str, max)
		if err != nil {
			t.Error(err)
		}

		resp, data, err := UnmarshalStringFromBytesData(data, max)
		if err != nil {
			t.Error(err)
		}

		if resp != str {
			t.Error("Unmarshal Fail")
		}

		if len(data) != 0 {
			t.Error("Unmarshal Return Data")
		}
		t.Log("ASD")

	}

	str := "123456"

	data, err := MarshalStringToBytes(str, 2)
	if err == nil {
		t.Error("Should error")
	}

	data, err = MarshalStringToBytes(str, 10)
	if err != nil {
		t.Error(err)
	}

	_, _, err = UnmarshalStringFromBytesData(data, 2)
	if err == nil {
		t.Error("should error")
	}

	b := &BoolHolder{true}
	data, err = b.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	a := new(BoolHolder)
	err = a.UnmarshalBinary(data)
	if err != nil {
		t.Error(err)
	}

	if a.Value != b.Value {
		t.Errorf("Should be same")
	}

}

func TestSanitizeFactomdLocation(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Normal", args{"http://courtesy-node.factom.com/"}, "http://courtesy-node.factom.com", false},
		{"SSL", args{"https://courtesy-node.factom.com/"}, "https://courtesy-node.factom.com", false},
		{"IP", args{"127.0.0.1"}, "http://127.0.0.1", false},
		{"IP&Port", args{"192.168.0.2:8088"}, "http://192.168.0.2:8088", false},
		{"Localhost", args{"localhost"}, "http://localhost", false},
		{"Localhost&Port", args{"localhost:80"}, "http://localhost:80", false},
		{"Short", args{"courtesy-node.factom.com"}, "http://courtesy-node.factom.com", false},
		{"With Port", args{"courtesy-node.factom.com:8088"}, "http://courtesy-node.factom.com:8088", false},
		{"With Path", args{"courtesy-node.factom.com/subpath"}, "http://courtesy-node.factom.com/subpath", false},
		{"With Trailing Slash", args{"courtesy-node.factom.com/subpath/"}, "http://courtesy-node.factom.com/subpath", false},
		{"With Path & Port", args{"courtesy-node.factom.com:443/subpath"}, "http://courtesy-node.factom.com:443/subpath", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeFactomdLocation(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeFactomdLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SanitizeFactomdLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}
