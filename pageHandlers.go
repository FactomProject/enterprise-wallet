package main

import (
	"fmt"
	"net/http"

	"github.com/FactomProject/factom"
)

// Every Handle struct must have settings
type SettingsStruct struct {
	Theme string // darkTheme or ""
}

var _ = fmt.Sprintf("")

func HandleIndexPage(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "indexPage", "")
	return nil
}

func HandleAddressBook(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "addressBook", "")
	return nil
}

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "settings", "")
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

	templates.ExecuteTemplate(w, "import-export-transaction", "")
	return nil
}

/*******************
 *  New Addresses  *
 *******************/

func HandleNewAddressFactoid(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "new-address-factoid", "")
	return nil
}

func HandleNewAddressEntryCredits(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "new-address-entry-credits", "")
	return nil
}

func HandleNewAddressExternal(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "new-address-external", "")
	return nil
}

/**************************
 *  Recieve/Send Factoids *
 **************************/

func HandleRecieveFactoids(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "receive-factoids", "")
	return nil
}

func HandleSendFactoids(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "send-factoids", "")
	return nil
}

/*******************
 *    Create EC    *
 *******************/

func HandleCreateEntryCredits(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "createEntryCredits", "")
	return nil
}

/*******************
 *    For Errors   *
 *******************/

func HandleNotFoundError(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	templates.ExecuteTemplate(w, "notFoundError", "")
	return nil
}
