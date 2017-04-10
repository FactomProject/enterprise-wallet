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
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/primitives/random"

	. "github.com/FactomProject/enterprise-wallet"
)

func ready() bool {
	r, err := factom.GetHeights()
	if err != nil || r.DirectoryBlockHeight < 1 {
		return false
	}
	return true
}

var TestWallet *wallet.WalletDB
var _ = fmt.Sprintf("")

func LoadTestWallet(port int) error {
	if TestWallet != nil { // If already instantiated
		return nil
	}

	wallet.GUI_DB = wallet.MAP
	wallet.WALLET_DB = wallet.MAP
	wallet.TX_DB = wallet.MAP

	wal, err := TestHelper.Start()
	if err != nil {
		return err
	}

	TestWallet = wal
	TestWallet.Wallet.TXDB().GetAllTXs()
	return nil
}

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

// It is big... It kept growing. Many variables are shared
// TODO: Break this up
func TestDisplayGETandPOST(t *testing.T) {
	for !ready() {
		time.Sleep(1 * time.Second)
	}
	var err error
	LoadTestWallet(7089)

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

	// POST Requests

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
				// Expected
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

		ecAdd := ""
		data, _ = handlePostRequestHelper("generate-new-address-ec", name)
		err = json.Unmarshal(data, resp)
		if err != nil {
			t.Errorf("Name is %s, err is: %s\n", name, err)
		} else if resp.Error != "none" {
			t.Error("Error response from request:", resp.Error)
		} else {
			ecAdd = resp.Content.Address
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
		if i%2 == 0 {
			data, _ = handlePostRequestHelper("display-private-key", `{"Address":"`+add+`"}`)
		} else {
			data, _ = handlePostRequestHelper("display-private-key", `{"Address":"`+ecAdd+`"}`)
		}
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
}

// Only add if not there
func importSandAddress() {
	_, _ = handlePostRequestHelper("new-address", `{"Name":"Sand","Secret":"Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"}`)
}

func TestSendEntryCreditsTransaction(t *testing.T) {
	for !ready() {
		time.Sleep(1 * time.Second)
	}
	// Import addresses with factoids
	// Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK - Sand
	// FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q

	// Es4CfH3dhjydTcUA5kD7maNWdDopi5kJBc3RkoVVPVv2HRAoLhC5
	// EC3JFMMpSpDEZFf7hBeSrcx25s6jkkoCV1F654J1uruBxNZRKCvF

	respA := new(jsonANPResponse)
	respG := new(jsonResponseGeneral)
	var err error
	importSandAddress()

	data, _ := handlePostRequestHelper("new-address", `{"Name":"Zero","Secret":"Es4CfH3dhjydTcUA5kD7maNWdDopi5kJBc3RkoVVPVv2HRAoLhC5"}`)
	err = json.Unmarshal(data, respA)
	if err != nil {
		t.Error("Failed importing address")
	}

	TestWallet.AddBalancesToAddresses()
	var currAmt int = 0
	data, _ = handlePostRequestHelper("get-address", `{"Address":"EC3JFMMpSpDEZFf7hBeSrcx25s6jkkoCV1F654J1uruBxNZRKCvF"}`)
	err = json.Unmarshal(data, respA)
	if err != nil || respA.Error != "none" {
		t.Error("Error occured getting address")
	} else {
		currAmt = int(respA.Content.Balance)
	}

	var totalSent int = 0
	for i := 0; i < 1; i++ {
		type jsonResponseRTS struct {
			Error   string            `json:"Error"`
			Content ReturnTransStruct `json:"Content"`
		}

		respR := new(jsonResponseRTS)

		sts := new(SendTransStruct)
		sts.TransType = "factoid"
		amt := random.RandIntBetween(5, 250)
		totalSent += amt
		sts.ToAmounts = []string{fmt.Sprintf("%d", amt)}
		sts.ToAddresses = []string{"EC3JFMMpSpDEZFf7hBeSrcx25s6jkkoCV1F654J1uruBxNZRKCvF"}

		data, err = json.Marshal(sts)
		if err != nil {
			t.Error(err)
		} else {
			jsonToSend := string(data)
			data, _ = handlePostRequestHelper("make-transaction", jsonToSend)
			err = json.Unmarshal(data, respR)
			if err != nil || respR.Error != "none" {

				t.Errorf("Error occured making transaction, %s", respR)
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
	fail := true
	trys := 0
	// try 3 times for correct ammount, sometimes it takes a little longer
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		TestWallet.AddBalancesToAddresses()
		time.Sleep(1 * time.Second)

		// Verify it worked
		data, _ = handlePostRequestHelper("get-address", `{"Address":"EC3JFMMpSpDEZFf7hBeSrcx25s6jkkoCV1F654J1uruBxNZRKCvF"}`)
		err = json.Unmarshal(data, respA)
		if err != nil || respA.Error != "none" {
			t.Error("Error occured getting address")
		} else {
			diff := (totalSent + currAmt) - int(respA.Content.Balance)
			if diff < 0 {
				diff = -1 * diff
			}

			if diff > 1 {
				trys++
			} else {
				fail = false
				break
			}
		}
	}

	if fail {
		t.Errorf("ECBuy: Tried %d times -- Balance is incorrect. Balance found is: %d, it should be %d\n CurrAmt: %d, TotalAdded: %d", trys, respA.Content.Balance, totalSent+currAmt, currAmt, totalSent)
	}
}

func TestConstructTransaction(t *testing.T) {
	for !ready() {
		time.Sleep(1 * time.Second)
	}
	// Import addresses with factoids
	// Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK - Sand
	// FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q

	// Fs25tDRiPrT9nmKpADm54ootJ5dk8yFmNAhmANsagVX85CFVg2GD - Zero
	// FA3WQuQigkTUH8jaA9EqtWd79pfULedXsje8qvF2ySZwLfNPQ8Ae
	respA := new(jsonANPResponse)
	respG := new(jsonResponseGeneral)
	var err error
	importSandAddress()

	data, _ := handlePostRequestHelper("new-address", `{"Name":"Zero","Secret":"Fs25tDRiPrT9nmKpADm54ootJ5dk8yFmNAhmANsagVX85CFVg2GD"}`)
	err = json.Unmarshal(data, respA)
	if err != nil {
		t.Error("Failed importing address")
	}

	TestWallet.AddBalancesToAddresses()
	var currAmt float64 = 0
	data, _ = handlePostRequestHelper("get-address", `{"Address":"FA3WQuQigkTUH8jaA9EqtWd79pfULedXsje8qvF2ySZwLfNPQ8Ae"}`)
	err = json.Unmarshal(data, respA)
	if err != nil || respA.Error != "none" {
		t.Error("Error occured getting address")
	} else {
		currAmt = float64(respA.Content.Balance) / 1e8
	}

	var totalSent float64 = 0
	for i := 0; i < 2; i++ {
		type jsonResponseRTS struct {
			Error   string            `json:"Error"`
			Content ReturnTransStruct `json:"Content"`
		}

		respR := new(jsonResponseRTS)

		sts := new(SendTransStruct)
		sts.TransType = "custom"
		if i == 1 {
			sts.TransType = "nosig"
		}
		amt := rand.Float64() * 10
		sts.ToAmounts = []string{fmt.Sprintf("%.8f", amt)}
		sts.ToAddresses = []string{"FA3WQuQigkTUH8jaA9EqtWd79pfULedXsje8qvF2ySZwLfNPQ8Ae"}

		sts.FromAddresses = []string{"FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q"}
		sts.FromAmounts = []string{fmt.Sprintf("%.8f", amt)}
		sts.FeeAddress = "FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q"

		sts.Signature = false

		data, err = json.Marshal(sts)
		if err != nil {
			t.Error(err)
		} else {
			jsonToSend := string(data)
			data, _ = handlePostRequestHelper("make-transaction", jsonToSend)
			err = json.Unmarshal(data, respR)
			if err != nil || respR.Error != "none" {
				t.Errorf("Error occured making transaction, %s", respR)
			} else {
				if sts.TransType == "nosig" {
					continue
				}
				totalSent += amt
				// lets send it
				data, _ = handlePostRequestHelper("send-transaction", jsonToSend)
				err = json.Unmarshal(data, respG)
				if err != nil || respG.Error != "none" {
					t.Errorf("Error occured sending transaction, %s", respG)
				} else {
				}
			}
		}
	}

	// Full block, blk times are 1 second in travis
	fail := true
	trys := 0
	// try 3 times for correct ammount, sometimes it takes a little longer
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		TestWallet.AddBalancesToAddresses()
		time.Sleep(1 * time.Second)

		// Verify it worked
		data, _ = handlePostRequestHelper("get-address", `{"Address":"FA3WQuQigkTUH8jaA9EqtWd79pfULedXsje8qvF2ySZwLfNPQ8Ae"}`)
		err = json.Unmarshal(data, respA)
		if err != nil || respA.Error != "none" {
			t.Error("Error occured getting address")
		} else {
			diff := (totalSent + currAmt) - (float64(respA.Content.Balance) / 1e8)
			if diff < 0 {
				diff = -1 * diff
			}

			if diff > 1 {
				trys++
			} else {
				fail = false
				break
			}
		}
	}

	if fail {
		t.Errorf("Construct:Tried %d times -- Balance is incorrect. Balance found is: %f, it should be %f\n CurrAmt: %f, TotalAdded: %f", trys, float64(respA.Content.Balance)/1e8, totalSent+currAmt, currAmt, totalSent)
	}
}

func TestSendFactoidsTransaction(t *testing.T) {
	for !ready() {
		time.Sleep(1 * time.Second)
	}
	// Import addresses with factoids
	// Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK - Sand
	// FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q

	// Fs2JQEA3DvhP7UFx7tCnrZZvfvnYkvD3eWwjs383PXuuHHXM8zph - Zero
	// FA2LsiAQTYKdYYxHLaBEhHsHDsnmpwayTyDzGRqQ8nAmsGwyLjRz
	respA := new(jsonANPResponse)
	respG := new(jsonResponseGeneral)
	var err error
	importSandAddress()

	data, _ := handlePostRequestHelper("new-address", `{"Name":"Zero","Secret":"Fs2JQEA3DvhP7UFx7tCnrZZvfvnYkvD3eWwjs383PXuuHHXM8zph"}`)
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
		currAmt = float64(respA.Content.Balance) / 1e8
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
		amt := rand.Float64() * 5
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
	fail := true
	trys := 0
	// try 3 times for correct ammount, sometimes it takes a little longer
	for i := 0; i < 3; i++ {
		time.Sleep(10 * time.Second)
		TestWallet.AddBalancesToAddresses()
		time.Sleep(1 * time.Second)

		// Verify it worked
		data, _ = handlePostRequestHelper("get-address", `{"Address":"FA2LsiAQTYKdYYxHLaBEhHsHDsnmpwayTyDzGRqQ8nAmsGwyLjRz"}`)
		err = json.Unmarshal(data, respA)
		if err != nil || respA.Error != "none" {
			t.Error("Error occured getting address")
		} else {
			diff := (totalSent + currAmt) - (float64(respA.Content.Balance) / 1e8)
			if diff < 0 {
				diff = -1 * diff
			}

			if diff > 1 {
				trys++
			} else {
				fail = false
				break
			}
		}
	}

	if fail {
		t.Errorf("FactoidSubmit:Tried %d times -- Balance is incorrect. Balance found is: %f, it should be %f\n CurrAmt: %f, TotalAdded: %f", trys, float64(respA.Content.Balance)/1e8, totalSent+currAmt, currAmt, totalSent)
	}
}

func TestValidAddresses(t *testing.T) {
	respS := new(jsonResponseGeneral)
	vectorValid := []string{
		"Fs1RM4pMUYV98mTZ2N2jKfT731bNSiNtiGdLaVZ7QhCYLuXZaGv4",
		"FA2Ax443J2xK63E38znfxrZ6kaVFmcfV7Hfy1Uj9dYV8LPB8m3Zf",
		"Fs1vPd2udNFy4c9zEpQHKu8pVpoQ7qFsGBHgxNgExVeynzoFXuFw",
		"FA1zs7rGN89Qf9CdjvMQT8sGChejobwvE97fr33VbhLDRFoHbXam",
		"Es4CfH3dhjydTcUA5kD7maNWdDopi5kJBc3RkoVVPVv2HRAoLhC5",
		"EC3JFMMpSpDEZFf7hBeSrcx25s6jkkoCV1F654J1uruBxNZRKCvF",
	}

	vectorInvalid := []string{
		"FA1RM4pMUYV98mTZ2N2jKfT731bNSiNtiGdLaVZ7QhCYLuXZaGv4",
		"FA3Ax443J2xK63E38znfxrZ6kaVFmcfV7Hfy1Uj9dYV8LPB8m3Zf",
		"Fs1vPd2udNFy4c9zEpQHKu8pVpoQ7qFsGBHgxNgExVeynzoFXuaw",
		"FA1zs7rGN89Qf9CdjvMQT8sGChejobwvE97fr33VbhLDRFHbXam",
		"EB4CfH3dhjydTcUA5kD7maNWdDopi5kJBc3RkoVVPVv2HRAoLhC5",
		"3JFMMpSpDEZFf7hBeSrcx25s6jkkoCV1F654J1uruBxNZRKCvF",
	}

	for _, a := range vectorValid {
		data, _ := handlePostRequestHelper("is-valid-address", a)
		err := json.Unmarshal(data, respS)
		if err != nil || respS.Error != "none" {
			t.Error("Failed on is-valid-address. Said a valid address is invalid")
		}
	}

	for _, a := range vectorInvalid {
		data, _ := handlePostRequestHelper("is-valid-address", a)
		err := json.Unmarshal(data, respS)
		if err != nil || respS.Content != "false" {
			t.Error("Failed on is-valid-address. Said an invalid address is valid")
		}
	}

	for i := 0; i < 100; i++ {
		data, _ := handlePostRequestHelper("is-valid-address", randomString(i))
		err := json.Unmarshal(data, respS)
		if err != nil || respS.Content != "false" {
			t.Error("Failed on is-valid-address. Said an invalid address is valid")
		}
	}

}

// 'yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow'
// FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1
func TestImportKoinify(t *testing.T) {
	respA := new(jsonANPResponse)
	k := "yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow"
	data, _ := handlePostRequestHelper("import-koinify", `{"Name":"Random","Koinify":"`+k+`"}`)
	err := json.Unmarshal(data, respA)
	if err != nil {
		t.Error(err)
	} else {
		if respA.Content.Address != "FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1" {
			t.Errorf("Koinify import failed. Expected FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1, got %s", respA.Content.Address)
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
