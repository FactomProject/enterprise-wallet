package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

var (
	FILES_PATH     string = "web/"
	templates      *template.Template
	MasterSettings *SettingsStruct

	mux           *http.ServeMux
	TemplateMutex sync.Mutex
)

// TODO: Compile statics into Go
func ServeWallet(port int) {
	templates = template.New("main")
	// Put function into templates
	funcMap := map[string]interface{}{"mkArray": mkArray, "compareInts": compareInts}
	templates.Funcs(template.FuncMap(funcMap))
	templates = template.Must(templates.ParseGlob(FILES_PATH + "templates/*.html"))

	// Start Settings
	// TODO: Load from DB
	MasterSettings = new(SettingsStruct)
	MasterSettings.Theme = ""

	// Mux for static files
	mux = http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web/statics")))

	http.HandleFunc("/", static(pageHandler))
	http.HandleFunc("/GET", HandleGETRequests)
	http.HandleFunc("/POST", HandlePOSTRequests)

	portStr := ":" + strconv.Itoa(port)

	fmt.Println("Starting Wallet on http://localhost" + portStr + "/")
	http.ListenAndServe(portStr, nil)
}

// Makes an array inside a template
func mkArray(args ...interface{}) []interface{} {
	return args
}

func compareInts(a int, b int) bool {
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

func pageHandler(w http.ResponseWriter, r *http.Request) {
	// Remove any GET data
	request := strings.Split(r.RequestURI, "?")
	fmt.Println(r.RequestURI)

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
	case "/receive-factoids":
		err = HandleRecieveFactoids(w, r)
	case "/send-factoids":
		err = HandleSendFactoids(w, r)
	default:
		err = HandleNotFoundError(w, r)
	}

	if err != nil {
		fmt.Printf("An error has occured")
	}
}

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

func jsonResp(content interface{}) []byte {
	e := newJsonResponse("none", content)
	return e.Bytes()
}

func jsonError(err string) []byte {
	e := newJsonResponse(err, "none")
	return e.Bytes()
	//return []byte("{'error':'" + err + "'}")
}

func HandleGETRequests(w http.ResponseWriter, r *http.Request) {
	// Only handles GET
	if r.Method != "GET" {
		return
	}
	req := r.FormValue("request")
	switch req {
	case "addresses":
		data, err := MasterWallet.GetGUIWalletJSON()
		if err != nil {
			w.Write(jsonError(err.Error()))
		}

		w.Write(data)
	default:
		w.Write(jsonError("Not a valid request"))
	}
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
		}
	case "display-private-key":
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

		_, list := MasterWallet.GetGUIAddress(a.Address)
		if list == 0 {
			w.Write(jsonError("Not found"))
			return
		}

		secret, err := MasterWallet.GetPrivateKey(a.Address)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		w.Write(jsonResp(secret))
	case "is-valid-address":
		add := r.FormValue("json")
		v := MasterWallet.IsValidAddress(add)
		if v {
			w.Write(jsonResp("true"))
		} else {
			w.Write(jsonResp("false"))
		}
	case "send-transaction":
		type SendTransStruct struct {
			Addresses []string `json:"OutputAddresses"`
			Amounts   []string `json:"OutputAmounts"`
		}

		trans := new(SendTransStruct)

		jsonElement := r.FormValue("json")
		err := json.Unmarshal([]byte(jsonElement), trans)
		if err != nil {
			w.Write(jsonError(err.Error()))
			return
		}

		name, err := MasterWallet.ConstructSendFactoidsStrings(trans.Addresses, trans.Amounts)
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

	default:
		w.Write(jsonError("Not a valid request"))
	}

}
