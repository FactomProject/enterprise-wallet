package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	//"strings"
	"math/rand"
	"testing"
	"time"

	"github.com/FactomProject/enterprise-wallet/TestHelper"
	"github.com/FactomProject/enterprise-wallet/address"
	"github.com/FactomProject/enterprise-wallet/wallet"
	"github.com/FactomProject/factomd/common/primitives/random"

	. "github.com/FactomProject/enterprise-wallet"
)

var TestWallet *wallet.WalletDB
var _ = fmt.Sprintf("")

func LoadTestWallet(port int) error {
	if TestWallet != nil { // If already instantiated
		return nil
	}

	wallet.GUI_DB = wallet.MAP
	wallet.WALLET_DB = wallet.MAP
	wallet.TX_DB = wallet.MAP

	wal, err := TestHelper.Start(port)
	if err != nil {
		return err
	}

	TestWallet = wal
	TestWallet.Wallet.TXDB().GetAllTXs()
	return nil
}

// It is big... It kept growing. Many variables are shared
// TODO: Break this up
func TestDisplayGETandPOST(t *testing.T) {
	var err error
	LoadTestWallet(7089)
	defer TestHelper.Stop()

	MasterWallet = TestWallet
	MasterSettings = new(SettingsStruct)
	InitTemplate()

	// TODO: Check outputs of GET requests

	// GET Requests
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8091/?request=synced", nil)
	HandleGETRequests(w, r)
	w = httptest.NewRecorder()

	r = httptest.NewRequest("GET", "localhost:8091/?request=addresses-no-bal", nil)
	HandleGETRequests(w, r)
	w = httptest.NewRecorder()

	r = httptest.NewRequest("GET", "localhost:8091/?request=addresses", nil)
	HandleGETRequests(w, r)
	w = httptest.NewRecorder()

	r = httptest.NewRequest("GET", "localhost:8091/?request=balances", nil)
	HandleGETRequests(w, r)
	w = httptest.NewRecorder()

	r = httptest.NewRequest("GET", "localhost:8091/?request=related-transactions", nil)
	HandleGETRequests(w, r)
	w = httptest.NewRecorder()

	r = httptest.NewRequest("GET", "localhost:8091/?request=on", nil)
	HandleGETRequests(w, r)
	w = httptest.NewRecorder()

	r = httptest.NewRequest("GET", "localhost:8091/?request=not-valid-random", nil)
	HandleGETRequests(w, r)
	w = httptest.NewRecorder()

	type jsonResponseGeneral struct {
		Error   string      `json:"Error"`
		Content interface{} `json:"Content"`
	}

	type jsonResponseStrings struct {
		Error   string `json:"Error"`
		Content string `json:"Content"`
	}

	// POST Requests
	type jsonANPResponse struct {
		Error   string                   `json:"Error"`
		Content *address.AddressNamePair `json:"Content"`
	}

	type SettingsToggle struct {
		Bools           []bool `json:"Values"` // A list of the boolean settings
		FactomdLocation string `json:"FactomdLocation"`
	}

	s := new(SettingsToggle)
	s.Bools = []bool{true, true, true, true}
	dataS, err := json.Marshal(s)
	if err != nil {
		t.Error(err)
		// Getting private keys will now fails
	} else {
		data, _ := handlePostRequestHelper("adjust-settings", string(dataS))
		respS := new(jsonResponseStrings)
		err = json.Unmarshal(data, respS)
		if err != nil {
			t.Error(err)
			// Getting private keys will now fails
		} else {
			if respS.Error != "none" {
				t.Error("Settings not updated:", respS.Error)
				// Getting private keys will now fails
			}
		}
	}
	respG := new(jsonResponseGeneral)
	resp := new(jsonANPResponse)
	// Make 5 addresses and change their names
	for i := 0; i < 30; i++ {
		name := randomString(i)
		// Generate Address
		var add string = ""
		data, _ := handlePostRequestHelper("generate-new-address-factoid", name)
		err = json.Unmarshal(data, resp)
		if err != nil {
			if i == 0 || i > 20 {
				// should not failt
				continue
			}
			t.Errorf("Name is %s, err is: %s\n", name, err)
		} else if resp.Error == "none" {
			if i == 0 || i > 20 {
				t.Error("This should fail. The name is too long or too short")
				continue
			}
			add = resp.Content.Address
		} else {
			t.Error("Error response from request:", resp.Error)
		}

		data, _ = handlePostRequestHelper("generate-new-address-ec", name)
		err = json.Unmarshal(data, resp)
		if err != nil {
			t.Errorf("Name is %s, err is: %s\n", name, err)
		} else if resp.Error != "none" {
			t.Error("Error response from request:", resp.Error)
		}

		// Change name
		respG = new(jsonResponseGeneral)
		if add != "" {
			data, _ := handlePostRequestHelper("address-name-change", `{"Name":"NewName", "Address":"`+add+`"}`)
			err = json.Unmarshal(data, respG)
			if err != nil {
				t.Error(err)
			} else if resp.Error != "none" {
				t.Error("Error response from request:", resp.Error)
			}

		} else {
			t.Error("Cannot test change name, generate address failed")
		}

		// Check name change && Get
		resp = new(jsonANPResponse)
		if add != "" {
			data, _ := handlePostRequestHelper("get-address", `{"Address":"`+add+`"}`)
			err = json.Unmarshal(data, resp)
			if err != nil {
				t.Error(err)
			} else if resp.Error != "none" {
				t.Error("Error response from request:", resp.Error)
			} else {
				if resp.Content.Name != "NewName" {
					t.Error("Name did not change")
				}
			}

		} else {
			t.Error("Cannot test change name success, generate address failed")
		}

		// Get a private key
		data, _ = handlePostRequestHelper("display-private-key", `{"Address":"`+add+`"}`)
		respS := new(jsonResponseStrings)
		err = json.Unmarshal(data, respS)
		if err != nil {
			t.Error(err)
		} else {
			// Got back a secret, lets check it
			data, _ = handlePostRequestHelper("new-address", `{"Name":"Taken", "Secret":"`+respS.Content+`"}`)
			respG := new(jsonResponseGeneral)
			err = json.Unmarshal(data, respG)
			if err != nil {
				t.Error(err)
			} else {
				if respG.Error != "Address already exists" {
					t.Error("This should error out, as the address already exists")
				}
			}

		}

	}

	// Make external addresses
	// Fs34PmX4gzBFDDwENuBAGGHeh7WLY6WJzVv6MiiYR8gYzCgsYsha
	// FA32vCmmaaB2ryHC35ZagviXgpvhMMuQ4tKD6m51Gg3nMW7UeNnK
	respA := new(jsonANPResponse)
	add := ""
	data, _ := handlePostRequestHelper("new-external-address", `{"Name":"Ext1", "Public":"FA32vCmmaaB2ryHC35ZagviXgpvhMMuQ4tKD6m51Gg3nMW7UeNnK"}`)
	err = json.Unmarshal(data, respA)
	if err != nil {
		t.Error(err)
	} else if respA.Error == "none" {
		add = respA.Content.Address
	} else {
		t.Error("Error response from request:", respA.Error)
	}

	// check if it failed or not
	if add == "" {
		t.Error("Add external address failed, cannot check delete external")
	} else {
		data, _ = handlePostRequestHelper("delete-address", `{"Name":"Ext1", "Address":"FA32vCmmaaB2ryHC35ZagviXgpvhMMuQ4tKD6m51Gg3nMW7UeNnK"}`)
		respG := new(jsonResponseGeneral)
		err = json.Unmarshal(data, respG)
		if err != nil {
			t.Error(err)
		} else {
			if respG.Error != "none" {
				t.Error("Error deleting external address:", err)
			}
		}
	}

	// Seed checks
	oldSeed := ""
	newSeed := ""
	respS := new(jsonResponseStrings)

	data, _ = handlePostRequestHelper("get-seed", "")
	err = json.Unmarshal(data, respS)
	if err != nil {
		t.Error(err)
	} else if respS.Error != "none" {
		t.Error(respS.Error)
	} else {
		oldSeed = respS.Content
	}

	respS = new(jsonResponseStrings)
	data, _ = handlePostRequestHelper("import-seed", `{"Seed":"shield hotel tent walk candy final smooth zebra island loan key hundred"}`)
	err = json.Unmarshal(data, respS)
	if err != nil {
		t.Error(err)
	} else if respS.Error != "none" {
		t.Error(respS.Error)
	} else {
		newSeed = respS.Content
	}

	if newSeed != "" {
		respS = new(jsonResponseStrings)
		data, _ = handlePostRequestHelper("get-seed", "")
		err = json.Unmarshal(data, respS)
		if err != nil {
			t.Error(err)
		} else {
			if oldSeed == respS.Content {
				t.Error("Seed was supposed to be changed.")
			}

			if newSeed != respS.Content {
				t.Error("Seed is unexpected value")
			}
		}
	} else {
		t.Error("Could not check seed change, import-seed failed")
	}

	// Import addresses with factoids
	// Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK - Sand
	// FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q

	// Fs2JQEA3DvhP7UFx7tCnrZZvfvnYkvD3eWwjs383PXuuHHXM8zph - Zero
	// FA2LsiAQTYKdYYxHLaBEhHsHDsnmpwayTyDzGRqQ8nAmsGwyLjRz
	data, _ = handlePostRequestHelper("new-address", `{"Name":"Sand","Secret":"Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"}`)
	err = json.Unmarshal(data, respA)
	if err != nil {
		t.Error("Failed importing address")
	}

	data, _ = handlePostRequestHelper("new-address", `{"Name":"Zero","Secret":"Fs2JQEA3DvhP7UFx7tCnrZZvfvnYkvD3eWwjs383PXuuHHXM8zph"}`)
	err = json.Unmarshal(data, respA)
	if err != nil {
		t.Error("Failed importing address")
	}

	TestWallet.AddBalancesToAddresses()
	var currAmt float64 = 0
	data, _ = handlePostRequestHelper("get-address", `{"Address":"FA2LsiAQTYKdYYxHLaBEhHsHDsnmpwayTyDzGRqQ8nAmsGwyLjRz"}`)
	err = json.Unmarshal(data, respA)
	if err != nil || respA.Error != "none" {
		t.Error("Error occured getting address")
	} else {
		currAmt = respA.Content.Balance
	}

	var totalSent float64 = 0
	for i := 0; i < 100; i++ {
		type jsonResponseRTS struct {
			Error   string            `json:"Error"`
			Content ReturnTransStruct `json:"Content"`
		}

		respR := new(jsonResponseRTS)

		sts := new(SendTransStruct)
		sts.TransType = "factoid"
		amt := rand.Float64() * 10
		totalSent += amt
		sts.ToAmounts = []string{fmt.Sprintf("%.8f", amt)}
		sts.ToAddresses = []string{"FA2LsiAQTYKdYYxHLaBEhHsHDsnmpwayTyDzGRqQ8nAmsGwyLjRz"}

		data, err = json.Marshal(sts)
		if err != nil {
			t.Error(err)
		} else {
			jsonToSend := string(data)
			data, _ = handlePostRequestHelper("make-transaction", jsonToSend)
			err = json.Unmarshal(data, respR)
			if err != nil || respR.Error != "none" {
				t.Error("Error occured making transaction")
			} else {
				// lets send it
				data, _ = handlePostRequestHelper("send-transaction", jsonToSend)
				err = json.Unmarshal(data, respG)
				if err != nil || respG.Error != "none" {
					t.Error("Error occured sending transaction")
				} else {
				}
			}
		}
	}

	// Full block, blk times are 1 second in travis
	time.Sleep(10 * time.Second)
	TestWallet.AddBalancesToAddresses()

	// Verify it worked
	data, _ = handlePostRequestHelper("get-address", `{"Address":"FA2LsiAQTYKdYYxHLaBEhHsHDsnmpwayTyDzGRqQ8nAmsGwyLjRz"}`)
	err = json.Unmarshal(data, respA)
	if err != nil || respA.Error != "none" {
		t.Error("Error occured getting address")
	} else {
		diff := (totalSent + currAmt) - respA.Content.Balance
		if diff < 0 {
			diff = -1 * diff
		}

		if diff > 1 {
			t.Errorf("Balance is incorrect. Balance found is: %f, it should be %f", respA.Content.Balance, totalSent+currAmt)
		}
	}
}

const StringAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

func randomString(l int) string {
	answer := []byte{}
	for i := 0; i < l; i++ {
		answer = append(answer, StringAlphabet[random.RandIntBetween(0, len(StringAlphabet)-1)])
	}
	return string(answer)
}

// handlePostRequestHelper returns the json result in bytes and a string
func handlePostRequestHelper(request string, json string) ([]byte, string) {
	form := url.Values{}
	w := httptest.NewRecorder()
	form.Add("request", request)
	form.Add("json", json)
	r := httptest.NewRequest("POST", "localhost:8091", nil)
	r.Form = form

	HandlePOSTRequests(w, r)

	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	return buf.Bytes(), buf.String()
}

func handleGETRequestHelper(request string) ([]byte, string) {
	form := url.Values{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8091/?request="+request, nil)
	r.Form = form

	HandleGETRequests(w, r)

	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	return buf.Bytes(), buf.String()
}
