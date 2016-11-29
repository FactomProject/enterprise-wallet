package wallet

/*
 * Manages all the addresses and 2 databases (Wallet DB and GUI DB)
 *
 */

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/FactomProject/M2WalletGUI/address"
	"github.com/FactomProject/M2WalletGUI/wallet/database"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
	// "github.com/FactomProject/factom/wallet/wsapi"
	"encoding/json"
	"github.com/FactomProject/factomd/common/interfaces"
)

const (
	MAP int = iota
	LDB
	BOLT
)

var (
	GUI_DB    = MAP
	WALLET_DB = MAP
)

// Wallet interacting with LDB and factom/wallet
//   The LDB doesn't need to be updated often, so we save after every add and only
//   deal with cached version
type WalletDB struct {
	GUIlDB        database.IDatabase        // GUI DB
	guiWallet     *WalletStruct             // Cached version on GUI LDB
	Wallet        *wallet.Wallet            // Wallet from factom/wallet
	TransactionDB *wallet.TXDatabaseOverlay // Used to display transactions

	// List of transactions related to any address in address book
	cachedTransactions []interfaces.ITransaction
	cachedHeight       int64
}

// For now is same as New
func LoadWalletDB() (*WalletDB, error) {
	return NewWalletDB()
}

func NewWalletDB() (*WalletDB, error) {
	w := new(WalletDB)

	// TODO: Adjust this path
	var db database.IDatabase
	var err error
	switch GUI_DB { // Decides type of wallet DB
	case MAP:
		db, err = database.NewMapDB()
	case LDB:
		db, err = database.NewLevelDB(GetHomeDir() + "/.factom/m2/gui_wallet_ldb")
	case BOLT:
	}
	if err != nil {
		return nil, err
	}

	w.GUIlDB = db

	w.guiWallet = NewWallet()
	data, err := w.GUIlDB.Get([]byte("wallet"))
	if err == nil {
		err = w.guiWallet.UnmarshalBinary(data)
	}
	if err != nil {
		data, err := w.guiWallet.MarshalBinary()
		if err != nil {
			return nil, err
		}

		err = w.GUIlDB.Put([]byte("wallet"), data)
		if err != nil {
			return nil, err
		}
	}

	// TODO: Adjust this path
	var wal *wallet.Wallet
	switch GUI_DB { // Decides type of wallet DB
	case MAP:
		wal, err = wallet.NewMapDBWallet()
	case LDB:
		wal, err = wallet.NewOrOpenLevelDBWallet(GetHomeDir() + "/.factom/m2/gui_wallet_testing")
	case BOLT:
		wal, err = wallet.NewOrOpenBoltDBWallet(GetHomeDir() + "/.factom/m2/gui_wallet_testing.db")
	}
	if err != nil {
		return nil, err
	}
	w.Wallet = wal

	txdb, err := wallet.NewTXBoltDB(fmt.Sprint(GetHomeDir(), "/.factom/wallet/factoid_blocks.cache"))
	if err != nil {
		return nil, fmt.Errorf("Could not add transaction database to wallet:", err)
	} else {
		w.Wallet.AddTXDB(txdb)
	}

	w.TransactionDB = w.Wallet.TXDB()

	err = w.UpdateGUIDB()
	if err != nil {
		return nil, err
	}

	// go wsapi.Start(w.Wallet, fmt.Sprintf(":%d", 8089), *(factom.RpcConfig))

	return w, nil
}

func (w *WalletDB) GetGUIWalletJSON() ([]byte, error) {
	w.AddBalancesToAddresses()
	return json.Marshal(w.guiWallet)
}

func (w *WalletDB) AddBalancesToAddresses() {
	w.guiWallet.AddBalancesToAddresses()
}

// Grabs the list of addresses from the walletDB and updates our GUI
// with any that are missing. All will be external
func (w *WalletDB) UpdateGUIDB() error {
	faAdds, ecAdds, err := w.Wallet.GetAllAddresses()
	if err != nil {
		return err
	}

	var names []string
	var addresses []string

	// Add addresses to GUI from cli
	for _, fa := range faAdds {
		_, list := w.GetGUIAddress(fa.String())
		if list == 0 {
			names = append(names, "FAImported-Undefined")
			addresses = append(addresses, fa.String())
		}
	}

	for _, ec := range ecAdds {
		_, list := w.GetGUIAddress(ec.String())
		if list == 0 {
			names = append(names, "ECImported-Undefined")
			addresses = append(addresses, ec.String())
		}
	}

	if len(names) > 0 {
		err = w.addBatchGUIAddresses(names, addresses)
		if err != nil {
			return err
		}
	}

	// Todo: Remove addresses that were deleted in cli?

	return w.Save()
}

func (w *WalletDB) Close() error {
	// Combine all close errors, as all need to get closed
	errCount := 0
	errString := ""

	err := w.Save()
	if err != nil {
		errCount++
		errString = errString + "; " + err.Error()
	}
	err = w.Wallet.Close()
	if err != nil {
		errCount++
		errString = errString + "; " + err.Error()
	}
	err = w.GUIlDB.Close()
	if err != nil {
		errCount++
		errString = errString + "; " + err.Error()
	}

	if errCount > 0 {
		return fmt.Errorf("Found %d errors: %s", errCount, errString)
	}
	return nil
}

func (w *WalletDB) Save() error {
	data, err := w.guiWallet.MarshalBinary()
	if err != nil {
		return err
	}

	err = w.GUIlDB.Put([]byte("wallet"), data)
	if err != nil {
		return err
	}

	return nil
}

func (w *WalletDB) GenerateFactoidAddress(name string) (*address.AddressNamePair, error) {
	address, err := w.Wallet.GenerateFCTAddress()

	if err != nil {
		return nil, err
	}

	anp, err := w.guiWallet.AddAddress(name, address.String(), 1)
	if err != nil {
		return nil, err
	}

	err = w.Save()
	if err != nil {
		return nil, err
	}
	return anp, nil
}

func (w *WalletDB) GetPrivateKey(address string) (secret string, err error) {
	if !factom.IsValidAddress(address) {
		return "", fmt.Errorf("Not a valid address")
	}

	if address[:2] == "FA" {
		return w.getFCTPrivateKey(address)
	} else if address[:2] == "EC" {
		return w.getECPrivateKey(address)
	}

	return "", fmt.Errorf("Not a public address")
}

func (w *WalletDB) getECPrivateKey(address string) (secret string, err error) {
	adds, err := w.Wallet.GetAllECAddresses()
	if err != nil {
		return "", err
	}

	for _, ec := range adds {
		if strings.Compare(ec.String(), address) == 0 {
			return ec.SecString(), nil
		}
	}

	return "", fmt.Errorf("Address not found")
}

func (w *WalletDB) getFCTPrivateKey(address string) (secret string, err error) {
	adds, err := w.Wallet.GetAllFCTAddresses()
	if err != nil {
		return "", err
	}

	for _, fa := range adds {
		if strings.Compare(fa.String(), address) == 0 {
			return fa.SecString(), nil
		}
	}

	return "", fmt.Errorf("Address not found")
}

func (w *WalletDB) GenerateEntryCreditAddress(name string) (*address.AddressNamePair, error) {
	address, err := w.Wallet.GenerateECAddress()
	if err != nil {
		return nil, err
	}

	anp, err := w.guiWallet.AddAddress(name, address.String(), 2)
	if err != nil {
		return nil, err
	}

	w.Save()
	if err != nil {
		return nil, err
	}

	return anp, nil
}

// TODO: Fix, make guiwallet take the remove
func (w *WalletDB) RemoveAddress(address string) (*address.AddressNamePair, error) {
	anp, err := w.guiWallet.RemoveAddress(address)
	if err != nil {
		return nil, err
	}

	err = w.Save()
	if err != nil {
		return nil, err
	}

	return anp, nil
}

func (w *WalletDB) AddAddress(name string, secret string) (*address.AddressNamePair, error) {
	if !factom.IsValidAddress(secret) {
		return nil, fmt.Errorf("Not a valid private key")
	} else if secret[:2] == "Fs" {
		add, err := factom.GetFactoidAddress(secret)
		if err != nil {
			return nil, err
		}

		err = w.Wallet.InsertFCTAddress(add)
		if err != nil {
			return nil, err
		}

		anp, err := w.addGUIAddress(name, add.String())
		if err != nil {
			return nil, err
		}

		err = w.Save()
		if err != nil {
			return nil, err
		}

		return anp, nil
	} else if secret[:2] == "Es" {
		add, err := factom.GetECAddress(secret)
		if err != nil {
			return nil, err
		}

		err = w.Wallet.InsertECAddress(add)
		if err != nil {
			return nil, err
		}

		anp, err := w.addGUIAddress(name, add.String())
		if err != nil {
			return nil, err
		}

		err = w.Save()
		if err != nil {
			return nil, err
		}

		return anp, nil
	}
	return nil, fmt.Errorf("Not a valid private key")
}

// Only adds to GUI Database
func (w *WalletDB) addBatchGUIAddresses(names []string, addresses []string) error {
	if len(names) != len(addresses) {
		return fmt.Errorf("List length does not match")
	}

	for i := 0; i < len(names); i++ {
		w.addGUIAddress(names[i], addresses[i])
	}

	return w.Save()
}

// Only adds to GUI database
func (w *WalletDB) addGUIAddress(name string, address string) (*address.AddressNamePair, error) {
	anp, err := w.guiWallet.AddAddress(name, address, 3)
	if err != nil {
		return nil, err
	}
	err = w.Save()
	if err != nil {
		return nil, err
	}

	return anp, nil
}

// Returns address with associated name
// List is 0 for not found, 1 for Factoid address, 2 for EC Address, 3 for External
func (w *WalletDB) GetGUIAddress(address string) (anp *address.AddressNamePair, list int) {
	anp, list, _ = w.guiWallet.GetAddress(address)
	return
}

func (w *WalletDB) ChangeAddressName(address string, toName string) error {
	err := w.guiWallet.ChangeAddressName(address, toName)
	if err != nil {
		return err
	}
	return w.Save()
}

func (w *WalletDB) GetTotalGUIAddresses() uint32 {
	return w.guiWallet.GetTotalAddressCount()
}

func (w *WalletDB) GetAllGUIAddresses() []address.AddressNamePair {
	return w.guiWallet.GetAllAddresses()
}

func (w *WalletDB) IsValidAddress(address string) bool {
	return factom.IsValidAddress(address)
}

func (w *WalletDB) GetECBalance() int64 {
	return w.guiWallet.ECTotal
}

func (w *WalletDB) GetFactoidBalance() int64 {
	return w.guiWallet.FactoidTotal
}

func GetHomeDir() string {
	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard
	// lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}
	return homeDir
}
