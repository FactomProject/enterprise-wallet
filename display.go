package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

var (
	FILES_PATH string = "web/"
	templates  *template.Template

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

	// Mux for static files
	mux = http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web/statics")))

	http.HandleFunc("/", static(pageHandler))
	http.HandleFunc("/GET", HandleGETRequests)

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
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	fmt.Println(r.RequestURI)

	switch r.RequestURI {
	case "/":
		templates.ExecuteTemplate(w, "indexPage", "")
	case "/AddressBook":
		templates.ExecuteTemplate(w, "addressBook", "")
	case "/Settings":
		templates.ExecuteTemplate(w, "settings", "")
	case "/create-entry-credits":
		templates.ExecuteTemplate(w, "createEntryCredits", "")
	case "/edit-address-entry-credits":
		templates.ExecuteTemplate(w, "edit-address-entry-credits", "")
	case "/edit-address-external":
		templates.ExecuteTemplate(w, "edit-address-external", "")
	case "/edit-address-factoid":
		templates.ExecuteTemplate(w, "edit-address-factoid", "")
	case "/import-export-transaction":
		templates.ExecuteTemplate(w, "import-export-transaction", "")
	case "/new-address-entry-credits":
		templates.ExecuteTemplate(w, "new-address-entry-credits", "")
	case "/new-address-external":
		templates.ExecuteTemplate(w, "new-address-external", "")
	case "/new-address-factoid":
		templates.ExecuteTemplate(w, "new-address-factoid", "")
	case "/receive-factoids":
		templates.ExecuteTemplate(w, "receive-factoids", "")
	case "/send-factoids":
		templates.ExecuteTemplate(w, "send-factoids", "")
	default:
		templates.ExecuteTemplate(w, "notFoundError", "")
	}
}

func jsonError(err string) []byte {
	return []byte("{'error':'" + err + "'}")
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
