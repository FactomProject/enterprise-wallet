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
	for i := 0; i < 1000; i++ {
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
}
