// Page Handlers
//
// Page handlers is the HTML serving of each page. All static files are compiled into the binary,
// meaning there is no need for absolute pathing. The static files follow a relative pathing scheme
// to mimic the normal non-compiled in behavior, which means you can turn off compilated statics when
// developing.
//		enterprise-wallet -compiled=false
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

const MAX_FACTOMDLOCATION_SIZE int = 30

// SettingsStruct
// Every Handle struct must have settings
// This is used on every page
type SettingsStruct struct {
	// Marshaled
	DarkTheme       bool
	KeyExport       bool // Allow export of private key
	CoinControl     bool
	ImportExport    bool //Transaction import/export
	FactomdLocation string

	// Not marshaled
	Theme            string // darkTheme or ""
	ControlPanelPort int
	Synced           bool
}

// Refresh refreshes the "synced" flag, and anything else that needs to be done
// before a page loads
func (s *SettingsStruct) Refresh() (leaderHeight int64, entryHeight int64, fblockHeight uint32) {
	var err error
	leaderHeight = 0
	entryHeight = 0
	fblockHeight = 0

	h, err := factom.GetHeights()
	if err != nil || h == nil {
		s.Synced = false
		return
	}

	leaderHeight = h.LeaderHeight
	entryHeight = h.EntryHeight

	fblockHeight, err = MasterWallet.Wallet.TXDB().FetchNextFBlockHeight()
	if err != nil {
		s.Synced = false
		return
	}

	// 1 block grace period
	if h != nil && (h.DirectoryBlockHeight >= (h.LeaderHeight - 1)) {
		if fblockHeight >= uint32(h.DirectoryBlockHeight) {
			s.Synced = true
			return
		}
	}
	s.Synced = false
	return
}

func (a *SettingsStruct) IsSameAs(b *SettingsStruct) bool {
	if a.DarkTheme != b.DarkTheme {
		return false
	}
	if a.KeyExport != b.KeyExport {
		return false
	}
	if a.CoinControl != b.CoinControl {
		return false
	}
	if a.ImportExport != b.ImportExport {
		return false
	}

	if a.Theme != b.Theme {
		return false
	}

	if a.FactomdLocation != b.FactomdLocation {
		return false
	}

	return true
}

func (s *SettingsStruct) SetFactomdLocation(factomdLocation string) {
	factom.SetFactomdServer(factomdLocation)
}

func (s *SettingsStruct) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	var b []byte

	b = strconv.AppendBool(b, s.DarkTheme)
	if s.DarkTheme {
		b = append(b, 0x00)
	}
	b = strconv.AppendBool(b, s.KeyExport)
	if s.KeyExport {
		b = append(b, 0x00)
	}
	b = strconv.AppendBool(b, s.CoinControl)
	if s.CoinControl {
		b = append(b, 0x00)
	}
	b = strconv.AppendBool(b, s.ImportExport)
	if s.ImportExport {
		b = append(b, 0x00)
	}

	buf.Write(b)

	data, err := MarshalStringToBytes(s.FactomdLocation, MAX_FACTOMDLOCATION_SIZE)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (s *SettingsStruct) UnmarshalBinary(data []byte) error {
	_, err := s.UnmarshalBinaryData(data)
	return err
}

func unmarshalBool(booldata []byte) (bool, error) {
	if booldata[4] == 0x00 {
		booldata = booldata[:4]
	}
	b, err := strconv.ParseBool(string(booldata))
	if err != nil {
		return false, err
	}

	return b, nil
}

func (s *SettingsStruct) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	newData = data

	s.DarkTheme, err = unmarshalBool(newData[:5])
	if err != nil {
		return data, err
	}
	newData = newData[5:]

	if s.DarkTheme {
		s.Theme = "darkTheme"
	} else {
		s.Theme = ""
	}

	s.KeyExport, err = unmarshalBool(newData[:5])
	if err != nil {
		return data, err
	}
	newData = newData[5:]

	s.CoinControl, err = unmarshalBool(newData[:5])
	if err != nil {
		return data, err
	}
	newData = newData[5:]

	s.ImportExport, err = unmarshalBool(newData[:5])
	if err != nil {
		return data, err
	}
	newData = newData[5:]

	switch {
	case len(newData) == 0: // v1 : No settings
		s.FactomdLocation = "localhost:8088" // Will be overwritten if changed anyhow
	case len(newData) == 30 && bytes.Compare(newData[28:30], []byte{0x00, 0x00}) == 0: // v2 : Settings at length 30
		//end := MAX_FACTOMDLOCATION_SIZE
		nameData := bytes.Trim(newData[:30], "\x00")
		s.FactomdLocation = fmt.Sprintf("%s", nameData)
		if s.FactomdLocation == "" {
			s.FactomdLocation = "localhost:8088"
		}
		newData = newData[30:]
	default: // current
		var loc string
		loc, newData, err = UnmarshalStringFromBytesData(newData, MAX_FACTOMDLOCATION_SIZE)
		if err != nil {
			return data, err
		}
		s.FactomdLocation = loc
	}

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

type HandleSettingsStruct struct {
	Settings *SettingsStruct

	Success bool
}

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	suc := r.FormValue("success")

	st := new(HandleSettingsStruct)
	st.Settings = MasterSettings

	st.Success = false
	if suc == "true" {
		st.Success = true
	}

	templates.ExecuteTemplate(w, "settings", st)
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
		reqStruct := fmt.Sprintf("%-10s: %s\n%10s: %s", "RequestUrl", r.URL, "RequestForm", r.Form)
		return nil, fmt.Errorf(" %s is an invalid address.\n%s", address, reqStruct)
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
		return fmt.Errorf("Error in HandleEditAddressFactoids(): %s", err.Error())
	}

	templates.ExecuteTemplate(w, "edit-address-factoid", e)
	return nil
}

func HandleEditAddressEntryCredits(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	e, err := handleEditAddress(w, r)
	if err != nil {
		return fmt.Errorf("Error in HandleEditAddressEntryCredits(): %s", err.Error())
	}

	templates.ExecuteTemplate(w, "edit-address-entry-credits", e)
	return nil
}

func HandleEditAddressExternal(w http.ResponseWriter, r *http.Request) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	e, err := handleEditAddress(w, r)
	if err != nil {
		return fmt.Errorf("Error in HandleEditAddressExternal(): %s", err.Error())
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

	if MasterSettings.ImportExport {
		templates.ExecuteTemplate(w, "import-export-transaction", NewPlaceHolderStruct())
	} else {
		templates.ExecuteTemplate(w, "notFoundError", NewPlaceHolderStruct())
	}
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

func HandleError(w http.ResponseWriter, r *http.Request, err error) error {
	TemplateMutex.Lock()
	defer TemplateMutex.Unlock()

	s := NewPlaceHolderStruct()
	s.Content = err.Error()
	templates.ExecuteTemplate(w, "customError", s)
	return nil
}
