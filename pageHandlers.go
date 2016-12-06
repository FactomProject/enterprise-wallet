package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FactomProject/factom"
)

type PlaceHolderStruct struct {
	Settings *SettingsStruct
	Content  interface{}
}

func NewPlaceHolderStruct() *PlaceHolderStruct {
	e := new(PlaceHolderStruct)
	e.Settings = MasterSettings
	e.Content = nil
	return e
}

// Every Handle struct must have settings
// This is used on every page
type SettingsStruct struct {
	// Marshaled
	DarkTheme bool
	KeyExport bool // Allow export of private key

	// Not marshaled
	Theme string // darkTheme or ""
}

func (s *SettingsStruct) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	var b []byte

	b = strconv.AppendBool(b, s.DarkTheme)
	b = strconv.AppendBool(b, s.KeyExport)

	buf.Write(b)

	return buf.Next(buf.Len()), nil
}

func (s *SettingsStruct) UnmarshalBinary(data []byte) error {
	_, err := s.UnmarshalBinaryData(data)
	return err
}

func (s *SettingsStruct) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	newData = data

	b, err := strconv.ParseBool(string(newData[:4]))
	if err != nil {
		return data, err
	}
	s.DarkTheme = b
	newData = newData[4:]

	if b {
		s.Theme = "darkTheme"
	}

	b, err = strconv.ParseBool(string(newData[:4]))
	if err != nil {
		return data, err
	}
	s.KeyExport = b
	newData = newData[4:]
	return
}

/*
func (s *SettingsStruct) FormatFactoid() {
	str := fmt.Sprintf("%f", s.FactoidBalance)
	arr := strings.Split(str, ".")
	s.FactoidFormatted = fmt.Sprintf("%s.<small>%s</small>", arr[0], arr[1])
}*/

var _ = fmt.Sprintf("")

func HandleIndexPage(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "indexPage", NewPlaceHolderStruct())
	return nil
}

func HandleAddressBook(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "addressBook", NewPlaceHolderStruct())
	return nil
}

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "settings", NewPlaceHolderStruct())
	return nil
}

/*******************
 *  Edit Addresses *
 *******************/

type EditAddressStruct struct {
	Settings *SettingsStruct

	Address string
	Name    string
}

func NewEditAddressStruct(address string, name string) *EditAddressStruct {
	e := new(EditAddressStruct)
	e.Settings = MasterSettings
	e.Address = address
	e.Name = name

	return e
}

func handleEditAddress(w http.ResponseWriter, r *http.Request) (*EditAddressStruct, error) {
	address := r.FormValue("address")

	if !factom.IsValidAddress(address) {
		return nil, fmt.Errorf("Invalid Address")
	}

	name := r.FormValue("name")

	e := NewEditAddressStruct(address, name)
	return e, nil

}

func HandleEditAddressFactoids(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	e, err := handleEditAddress(w, r)
	if err != nil {
		return err
	}

	templates.ExecuteTemplate(w, "edit-address-factoid", e)
	return nil
}

func HandleEditAddressEntryCredits(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	e, err := handleEditAddress(w, r)
	if err != nil {
		return err
	}

	templates.ExecuteTemplate(w, "edit-address-entry-credits", e)
	return nil
}

func HandleEditAddressExternal(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	e, err := handleEditAddress(w, r)
	if err != nil {
		return err
	}

	templates.ExecuteTemplate(w, "edit-address-external", e)
	return nil
}

/*******************
 *  Import/Export  *
 *******************/

func HandleImportExportTransaction(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "import-export-transaction", NewPlaceHolderStruct())
	return nil
}

/*******************
 *  New Addresses  *
 *******************/

func HandleNewAddress(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "new-address", NewPlaceHolderStruct())
	return nil
}

func HandleNewAddressFactoid(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "new-address-factoid", NewPlaceHolderStruct())
	return nil
}

func HandleNewAddressEntryCredits(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "new-address-entry-credits", NewPlaceHolderStruct())
	return nil
}

func HandleNewAddressExternal(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "new-address-external", NewPlaceHolderStruct())
	return nil
}

/**************************
 *  Receive/Send Factoids *
 **************************/

type ReceiveFactoidsStruct struct {
	Settings *SettingsStruct

	Address string
	Name    string
}

func HandleReceiveFactoids(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	address := r.FormValue("address")
	name := r.FormValue("name")

	st := new(ReceiveFactoidsStruct)
	st.Settings = MasterSettings

	st.Address = address
	st.Name = name
	if MasterWallet.IsValidAddress(address) {
		templates.ExecuteTemplate(w, "receive-factoids", st)
	} else {
		templates.ExecuteTemplate(w, "receive-factoids", st)
	}
	return nil
}

func HandleSendFactoids(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "send-factoids", NewPlaceHolderStruct())
	return nil
}

/*******************
 *    Create EC    *
 *******************/

func HandleCreateEntryCredits(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "createEntryCredits", NewPlaceHolderStruct())
	return nil
}

/*******************
 *    For Errors   *
 *******************/

func HandleNotFoundError(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "notFoundError", NewPlaceHolderStruct())
	return nil
}
