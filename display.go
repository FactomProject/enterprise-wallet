package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/FactomProject/enterprise-wallet/web/files"
)

var (
	FILES_PATH     string = "web/"
	templates      *template.Template
	MasterSettings *SettingsStruct

	mux           *http.ServeMux
	TemplateMutex sync.Mutex

	// How often to update balances in cache
	BALANCE_UPDATE_INTERVAL time.Duration = 10 * time.Second
)

// Use or no use compiled statics. Keeping a non-compiled
// option for front end design changes
var COMPILED_STATICS = true

func SaveSettings() error {
	err := MasterWallet.GUIlDB.Put([]byte("gui-wallet"), []byte("settings"), MasterSettings)
	if err != nil {
		return err
	}
	err = MasterWallet.GUIlDB.Put([]byte("gui-wallet"), []byte("backed-up"), &BoolHolder{MasterSettings.BackedUp})
	return err
}

func InitTemplate() {
	TemplateMutex.Lock()
	// Put function into templates
	funcMap := map[string]interface{}{"mkArray": mkArray, "compareInts": compareInts, "compareStrings": compareStrings}
	templates = template.New("main")
	templates.Funcs(template.FuncMap(funcMap))
	if COMPILED_STATICS { // Use compiled
		templates = files.CustomParseGlob(templates, "templates/*.html")
		templates = template.Must(templates, nil)
	} else { // Use non-compiled
		fmt.Println(FILES_PATH)
		templates = template.Must(templates.ParseGlob(FILES_PATH + "templates/*.html"))
	}
	templates.Funcs(template.FuncMap(funcMap))
	TemplateMutex.Unlock()
}

func ServeWallet(port int) {
	// Templates
	InitTemplate()

	// Update the balances every 10 seconds to keep it updated. We can force
	// an update if we send a transaction or something
	go doEvery(BALANCE_UPDATE_INTERVAL, updateBalances)

	// Load the initial transaction DB. This takes some time, should start before user hits first page
	go MasterWallet.GetRelatedTransactions()

	// Mux for static files
	mux = http.NewServeMux()
	if COMPILED_STATICS {
		mux.Handle("/", files.StaticServer)
	} else {
		mux.Handle("/", http.FileServer(http.Dir(FILES_PATH+"statics")))
	}

	http.HandleFunc("/", static(pageHandler))
	http.HandleFunc("/GET", HandleGETRequests)
	http.HandleFunc("/POST", HandlePOSTRequests)

	portStr := "localhost:" + strconv.Itoa(port)

	fmt.Println("Starting GUI on http://" + portStr + "/")
	http.ListenAndServe(portStr, nil)
}

// mkArray makes an array inside a template
func mkArray(args ...interface{}) []interface{} {
	return args
}

// compareInts is used inside templates to compare ints
func compareInts(a int, b int) bool {
	return (a == b)
}

// compareStrings used inside templates to compare strings
func compareStrings(a string, b string) bool {
	return (a == b)
}

// For all static files. (CSS, JS, IMG, etc...)
func static(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ContainsRune(r.URL.Path, '.') {
			mux.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
}

// updateBalances updates various elements. Faster load times for user if these
// are loaded when they are not asking
func updateBalances(time.Time) {
	MasterWallet.AddBalancesToAddresses()
	MasterWallet.UpdateGUIDB()
	MasterWallet.GetRelatedTransactions()
}

// doEvery
// For go routines. Calls function once each duration.
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

// pageHandler redirects all page requests to proper handlers
func pageHandler(w http.ResponseWriter, r *http.Request) {
	MasterSettings.Synced = false
	request := strings.Split(r.RequestURI, "?")
	var err error
	switch request[0] {
	case "/":
		err = HandleIndexPage(w, r)
	case "/AddressBook":
		err = HandleAddressBook(w, r)
	case "/Settings":
		err = HandleSettings(w, r)
	case "/create-entry-credits":
		err = HandleCreateEntryCredits(w, r)
	case "/edit-address-entry-credits":
		err = HandleEditAddressEntryCredits(w, r)
	case "/edit-address-external":
		err = HandleEditAddressExternal(w, r)
	case "/edit-address-factoid":
		err = HandleEditAddressFactoids(w, r)
	case "/import-export-transaction":
		err = HandleImportExportTransaction(w, r)
	case "/new-address-entry-credits":
		err = HandleNewAddressEntryCredits(w, r)
	case "/new-address-external":
		err = HandleNewAddressExternal(w, r)
	case "/new-address-factoid":
		err = HandleNewAddressFactoid(w, r)
	case "/new-address":
		err = HandleNewAddress(w, r)
	case "/receive-factoids":
		err = HandleReceiveFactoids(w, r)
	case "/send-factoids":
		err = HandleSendFactoids(w, r)
	case "/Backup":
		err = HandleBackup0(w, r)
	case "/backup1":
		err = HandleBackup1(w, r)
	case "/backup2":
		err = HandleBackup2(w, r)
	case "/backup3":
		err = HandleBackup3(w, r)
	case "/import-seed":
		err = HandleImportSeed(w, r)
	case "/success-screen-backup":
		err = HandleSuccessScreen(w, r)
	case "/success-screen-import":
		err = HandleSuccessScreen(w, r)
	default:
		err = HandleNotFoundError(w, r)
	}

	if err != nil {
		fmt.Printf("An error has occurred")
		HandleError(w, r, err)
	}
}

// jsonResponse is used for responding to Post/Get Requests
type jsonResponse struct {
	Error   string      `json:"Error"`
	Content interface{} `json:"Content"`
}

func newJsonResponse(err string, content interface{}) *jsonResponse {
	j := new(jsonResponse)
	j.Error = err
	j.Content = content

	return j
}

func (j *jsonResponse) Bytes() []byte {
	data, err := json.Marshal(j)
	if err != nil {
		return nil
	}

	return data
}

// jsonResp used if request is successful
func jsonResp(content interface{}) []byte {
	e := newJsonResponse("none", content)
	return e.Bytes()
}

// jsonError used if request has an error
func jsonError(err string) []byte {
	e := newJsonResponse(err, "none")
	return e.Bytes()
}

func HandleGETRequests(w http.ResponseWriter, r *http.Request) {
	// Only handles GET
	if r.Method != "GET" {
		return
	}
	req := r.FormValue("request")
	switch req {
	case "on":
		w.Write(jsonResp(true))
	case "synced":
		type SyncedStruct struct {
			Synced       bool
			LeaderHeight int64
			EntryHeight  int64
			FblockHeight uint32
			Stage        int
		}
		s := new(SyncedStruct)

		lh, eh, fh := MasterSettings.Refresh()
		s.Synced = MasterSettings.Synced
		s.LeaderHeight = lh
		s.EntryHeight = eh
		s.FblockHeight = fh
		s.Stage = MasterWallet.GetStage()
		w.Write(jsonResp(s))
	case "addresses-no-bal":
		data, err := MasterWallet.GetGUIWalletJSON(false)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		w.Write(data)
	case "addresses":
		data, err := MasterWallet.GetGUIWalletJSON(true)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		w.Write(data)
	case "balances":
		bals := struct {
			EC int64
			FC int64
		}{MasterWallet.GetECBalance(), MasterWallet.GetFactoidBalance()}
		data := jsonResp(bals)
		if data != nil {
			w.Write(data)
			return
		}

		w.Write(jsonError("Error occurred"))
	case "related-transactions":
		if on, server := MasterWallet.FactomdOnline(); !on {
			errorMsg := fmt.Sprintf("Unable to connect to the factomd instance. Make sure there is an active, fully synced factomd at %s", server)
			w.Write(jsonError(errorMsg))
			return
		}

		trans, err := MasterWallet.GetRelatedTransactions()
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		} else {
			MasterWallet.ActiveCachedTransactions = trans
			if len(trans) > 100 {
				next := trans[:100]
				next = MasterWallet.ScrubDisplayTransactionsForNameChanges(next)
				w.Write(jsonResp(next))
			} else {
				next := trans
				next = MasterWallet.ScrubDisplayTransactionsForNameChanges(next)
				w.Write(jsonResp(next))
			}
		}
	default:
		w.Write(jsonError("Not a valid request"))
	}
}

// SendTransStruct is a struct for sending transactions
type SendTransStruct struct {
	TransType   string   `json:"TransType"`
	ToAddresses []string `json:"OutputAddresses"`
	ToAmounts   []string `json:"OutputAmounts"`

	FromAddresses []string `json:"InputAddresses"`
	FromAmounts   []string `json:"InputAmounts"`
	FeeAddress    string   `json:"FeeAddress"`

	Signature bool `json:"Signature, omitempty"`
}

type ReturnTransStruct struct {
	Name  string `json:"Name"`
	Total uint64 `json:"Total"`
	Fee   uint64 `json:"Fee"`
	Json  string `json:"Json"`
}

func HandlePOSTRequests(w http.ResponseWriter, r *http.Request) {
	// Only handles POST
	if r.Method != "POST" {
		return
	}

	// Form:
	//	request -- Request Function
	//	json	-- json object

	req := r.FormValue("request")
	switch req {
	case "address-name-change":
		type ANC struct {
			Address string `json:"Address"`
			ToName  string `json:"Name"`
		}
		j := r.FormValue("json")
		anc := new(ANC)
		err := json.Unmarshal([]byte(j), anc)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}
		err = MasterWallet.ChangeAddressName(anc.Address, anc.ToName)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		} else {
			w.Write(jsonResp("Success"))
		}
	case "delete-address":
		type ANC struct {
			Address string `json:"Address"`
			Name    string `json:"Name"`
		}
		j := r.FormValue("json")
		anc := new(ANC)
		err := json.Unmarshal([]byte(j), anc)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		_, list := MasterWallet.GetGUIAddress(anc.Address)
		if list != 3 {
			w.Write(jsonError("You can only delete External Addresses."))
			return
		}

		_, err = MasterWallet.RemoveAddress(anc.Address, list)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		} else {
			w.Write(jsonResp("Success"))
		}
	case "display-private-key":
		type Add struct {
			Address string `json:"Address"`
		}

		if !MasterSettings.KeyExport {
			w.Write(jsonResp("Displaying private key disabled in settings"))
			return
		}

		j := r.FormValue("json")
		a := new(Add)
		err := json.Unmarshal([]byte(j), a)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		_, list := MasterWallet.GetGUIAddress(a.Address)
		if list == -1 {
			w.Write(jsonError("Not found"))
			return
		}

		secret, err := MasterWallet.GetPrivateKey(a.Address)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		w.Write(jsonResp(secret))
	case "get-address":
		type Add struct {
			Address string `json:"Address"`
		}

		j := r.FormValue("json")
		a := new(Add)
		err := json.Unmarshal([]byte(j), a)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		anp, list := MasterWallet.GetGUIAddress(a.Address)
		if list == -1 {
			w.Write(jsonError("Not found"))
			return
		}

		w.Write(jsonResp(anp))
	case "is-valid-address":
		add := r.FormValue("json")
		v := MasterWallet.IsValidAddress(add)
		if v {
			w.Write(jsonResp("true"))
		} else {
			w.Write(jsonResp("false"))
		}
	case "generate-new-address-factoid":
		name := r.FormValue("json")
		anp, err := MasterWallet.GenerateFactoidAddress(name)
		if err != nil {
			w.Write(jsonError(err.Error()))
		} else {
			w.Write(jsonResp(anp))
		}
	case "generate-new-address-ec":
		name := r.FormValue("json")
		anp, err := MasterWallet.GenerateEntryCreditAddress(name)
		if err != nil {
			w.Write(jsonError(err.Error()))
		} else {
			w.Write(jsonResp(anp))
		}
	case "new-address":
		type NewAddressStruct struct {
			Name   string `json:"Name"`
			Secret string `json:"Secret"`
		}

		nas := new(NewAddressStruct)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), nas)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		anp, err := MasterWallet.AddAddress(nas.Name, nas.Secret)
		if err != nil {
			w.Write(jsonError(err.Error()))
		} else {
			w.Write(jsonResp(anp))
		}
	case "import-koinify":
		type NewKoinifyStruct struct {
			Name    string `json:"Name"`
			Koinify string `json:"Koinify"`
		}

		nas := new(NewKoinifyStruct)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), nas)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		anp, err := MasterWallet.ImportKoinify(nas.Name, nas.Koinify)
		if err != nil {
			w.Write(jsonError(err.Error()))
		} else {
			w.Write(jsonResp(anp))
		}
	case "new-external-address":
		type NewAddressStruct struct {
			Name   string `json:"Name"`
			Public string `json:"Public"`
		}

		nas := new(NewAddressStruct)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), nas)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		anp, err := MasterWallet.AddExternalAddress(nas.Name, nas.Public)
		if err != nil {
			w.Write(jsonError(err.Error()))
		} else {
			w.Write(jsonResp(anp))
		}
	case "get-needed-input":
		trans := new(SendTransStruct)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), trans)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		needed, err := MasterWallet.CalculateNeededInput(trans.ToAddresses, trans.ToAmounts)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		w.Write(jsonResp(needed))
	case "import-transaction":
		// new(SendTransStruct)
		transHex := r.FormValue("json")
		MasterWallet.Wallet.DeleteTransaction("importedTX")
		err := MasterWallet.ImportTransaction("importedTX", transHex)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		trans := MasterWallet.Wallet.GetTransactions()["importedTX"]
		if trans == nil {
			w.Write(jsonError("Transaction had an error importing."))
			return
		}

		transRet := new(SendTransStruct)
		inputs := trans.GetInputs()
		for _, in := range inputs {
			transRet.FromAddresses = append(transRet.FromAddresses, MasterWallet.FactoidAddressToHumanReadable(in.GetAddress()))
			transRet.FromAmounts = append(transRet.FromAmounts, fmt.Sprintf("%f", float64(in.GetAmount())/1e8))
		}

		outputs := trans.GetOutputs()
		for _, out := range outputs {
			transRet.ToAddresses = append(transRet.ToAddresses, MasterWallet.FactoidAddressToHumanReadable(out.GetAddress()))
			transRet.ToAmounts = append(transRet.ToAmounts, fmt.Sprintf("%f", float64(out.GetAmount())/1e8))
		}

		ecouts := trans.GetECOutputs()
		for _, out := range ecouts {
			transRet.ToAddresses = append(transRet.ToAddresses, MasterWallet.ECAddressToHumanReadable(out.GetAddress()))
			transRet.ToAmounts = append(transRet.ToAmounts, fmt.Sprintf("%d", out.GetAmount()))
		}

		err = trans.ValidateSignatures()
		if err == nil {
			transRet.Signature = true
		} else {
			transRet.Signature = false
		}

		w.Write(jsonResp(transRet))
	case "broadcast-transaction":
		err := MasterWallet.Wallet.SignTransaction("importedTX", true)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		txid, err := MasterWallet.SendTransaction("importedTX")
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		w.Write(jsonResp(txid))

	case "make-transaction":
		trans := new(SendTransStruct)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), trans)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		var r ReturnTransStruct

		name := ""
		if trans.TransType == "factoid" {
			newName, rt, err := MasterWallet.ConstructSendFactoidsStrings(trans.ToAddresses, trans.ToAmounts)
			if err != nil {
				MasterWallet.DeleteTransaction(name)
				w.Write(jsonError(err.Error()))
				return
			}

			name = newName
			r.Total = rt.Total
			r.Fee = rt.Fee
		} else if trans.TransType == "ec" {
			newName, rt, err := MasterWallet.ConstructConvertEntryCreditsStrings(trans.ToAddresses, trans.ToAmounts)
			if err != nil {
				MasterWallet.DeleteTransaction(name)
				w.Write(jsonError(err.Error()))
				return
			}

			name = newName
			r.Total = rt.Total
			r.Fee = rt.Fee
		} else if trans.TransType == "custom" {
			newName, rt, err := MasterWallet.ConstructTransactionFromValuesStrings(
				trans.ToAddresses, trans.ToAmounts, trans.FromAddresses, trans.FromAmounts, trans.FeeAddress, true)
			if err != nil {
				MasterWallet.DeleteTransaction(name)
				w.Write(jsonError(err.Error()))
				return
			}

			name = newName
			r.Total = rt.Total
			r.Fee = rt.Fee
		} else if trans.TransType == "nosig" {
			newName, rt, err := MasterWallet.ConstructTransactionFromValuesStrings(
				trans.ToAddresses, trans.ToAmounts, trans.FromAddresses, trans.FromAmounts, trans.FeeAddress, false)
			if err != nil {
				MasterWallet.DeleteTransaction(name)
				w.Write(jsonError(err.Error()))
				return
			}

			name = newName
			r.Total = rt.Total
			r.Fee = rt.Fee
		} else {
			w.Write(jsonError("Not a valid type"))
			return
		}

		j, err := MasterWallet.ExportTransaction(name)
		if err != nil {
			r.Json = "Error exporting transaction."
		} else {
			r.Json = j
		}

		r.Name = name
		w.Write(jsonResp(r))
	case "send-transaction":
		MasterWallet.BalanceCache.Clear()

		trans := new(SendTransStruct)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), trans)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		name, err := MasterWallet.CheckTransactionAndGetName(trans.ToAddresses, trans.ToAmounts, trans.FeeAddress)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		tHash, err := MasterWallet.SendTransaction(name)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		w.Write(jsonResp(tHash))
	case "adjust-settings":
		type SettingsToggle struct {
			Bools           []bool `json:"Values"` // A list of the boolean settings
			FactomdLocation string `json:"FactomdLocation"`
		}

		st := new(SettingsToggle)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), st)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		MasterSettings.DarkTheme = st.Bools[0]
		if st.Bools[0] {
			MasterSettings.Theme = "darkTheme"
		} else {
			MasterSettings.Theme = ""
		}
		MasterSettings.KeyExport = st.Bools[1]
		MasterSettings.CoinControl = st.Bools[2]
		MasterSettings.ImportExport = st.Bools[3]

		if loc, err := SanitizeFactomdLocation(st.FactomdLocation); err != nil {
			w.Write(jsonError(err.Error()))
			return
		} else {
			st.FactomdLocation = loc
		}

		fdChange := false
		if len(st.FactomdLocation) > 0 && st.FactomdLocation != MasterSettings.FactomdLocation {
			MasterSettings.FactomdLocation = st.FactomdLocation
			MasterSettings.SetFactomdLocation(MasterSettings.FactomdLocation)
			fdChange = true
		}

		err = SaveSettings()
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}
		if fdChange {
			w.Write(jsonResp(fmt.Sprintf("Settings updated and the location of factomd was changed to %s\n", MasterSettings.FactomdLocation)))
		} else {
			w.Write(jsonResp("Settings updated"))
		}
	case "get-seed":
		seed, err := MasterWallet.ExportSeed()
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}
		w.Write(jsonResp(seed))
	case "import-seed":
		type SeedStruct struct {
			Seed string `json:"Seed"`
		}

		ss := new(SeedStruct)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), ss)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		err = MasterWallet.ImportSeed(ss.Seed)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}
		w.Write(jsonResp(ss.Seed))
	case "more-cached-transaction":
		type MoreRelatedTransactionReq struct {
			Current int `json:"Current"` // Current index in list
			More    int `json:"More"`    // How many more
		}

		rt := new(MoreRelatedTransactionReq)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), rt)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		total := len(MasterWallet.ActiveCachedTransactions)
		max := rt.Current + rt.More
		if max > total {
			if rt.Current >= total {
				w.Write(jsonResp(nil))
				return
			}
			next := MasterWallet.ActiveCachedTransactions[rt.Current:]
			next = MasterWallet.ScrubDisplayTransactionsForNameChanges(next)
			w.Write(jsonResp(next))
		} else {
			next := MasterWallet.ActiveCachedTransactions[rt.Current:max]
			next = MasterWallet.ScrubDisplayTransactionsForNameChanges(next)
			w.Write(jsonResp(next))
		}

	default:
		w.Write(jsonError("Not a post valid request"))
	}

}
