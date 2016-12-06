package wallet

/*
 * Manages all the addresses and 2 databases (Wallet DB and GUI DB)
 *
 */

import (
	"fmt"
	"os"
	"os/user"
	"sort"
	"strings"

	"github.com/FactomProject/M2GUIWallet/address"
	"github.com/FactomProject/M2GUIWallet/wallet/database"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factomd/common/primitives"
	// "github.com/FactomProject/factom/wallet/wsapi"
	"encoding/json"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/database/mapdb"
)

const (
	MAP int = iota
	LDB
	BOLT
)

var (
	GUI_DB    = MAP
	WALLET_DB = MAP
	TX_DB     = MAP
)

// Wallet interacting with LDB and factom/wallet
//   The LDB doesn't need to be updated often, so we save after every add and only
//   deal with cached version
type WalletDB struct {
	GUIlDB        interfaces.IDatabase      //database.IDatabase        // GUI DB
	guiWallet     *WalletStruct             // Cached version on GUI LDB
	Wallet        *wallet.Wallet            // Wallet from factom/wallet
	TransactionDB *wallet.TXDatabaseOverlay // Used to display transactions

	// List of transactions related to any address in address book
	cachedTransactions []DisplayTransaction
	cachedHeight       uint32
}

// For now is same as New
func LoadWalletDB() (*WalletDB, error) {
	return NewWalletDB()
}

func NewWalletDB() (*WalletDB, error) {
	w := new(WalletDB)

	// TODO: Adjust this path
	var db interfaces.IDatabase
	var err error
	switch GUI_DB { // Decides type of wallet DB
	case MAP:
		db, err = database.NewMapDB()
	case LDB:
		db, err = database.NewOrOpenLevelDBWallet(GetHomeDir() + "/.factom/m2/gui_wallet_ldb")
	case BOLT:
		db, err = database.NewOrOpenBoltDBWallet(GetHomeDir() + "/.factom/m2/gui_wallet_bolt")
	}
	if err != nil {
		return nil, err
	}

	w.GUIlDB = db

	// Adds Wallet
	w.guiWallet = NewWallet()
	data, err := w.GUIlDB.Get([]byte("gui-wallet"), []byte("wallet"), new(WalletStruct))
	if err != nil || data == nil {
		err = w.GUIlDB.Put([]byte("gui-wallet"), []byte("wallet"), w.guiWallet)
		if err != nil {
			return nil, err
		}
	} else {
		w.guiWallet = data.(*WalletStruct)
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

	// TODO: Adjust this path
	var txdb *wallet.TXDatabaseOverlay
	switch TX_DB {
	case MAP:
		txdb = wallet.NewTXOverlay(new(mapdb.MapDB))
		err = nil
	case LDB:
		txdb, err = wallet.NewTXLevelDB(fmt.Sprint(GetHomeDir(), "/.factom/m2/wallet/factoid_blocks_level"))
	case BOLT:
		txdb, err = wallet.NewTXBoltDB(fmt.Sprint(GetHomeDir(), "/.factom/m2/wallet/factoid_blocks.cache"))
	}

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

	//w.cachedTransactions = make(interfaces.ITransaction)
	w.cachedHeight = 0

	// go wsapi.Start(w.Wallet, fmt.Sprintf(":%d", 8089), *(factom.RpcConfig))

	return w, nil
}

type TransactionAddressInfo struct {
	Name    string
	Address string
	Amount  uint64
	Type    string // FCT or EC
}

func NewTransactionAddressInfo(name string, address string, amount uint64, tokenType string) *TransactionAddressInfo {
	t := new(TransactionAddressInfo)
	t.Name = name
	t.Address = address
	t.Amount = amount
	t.Type = tokenType

	return t
}

// Names are "" if not in wallet
type DisplayTransaction struct {
	Inputs     []TransactionAddressInfo
	TotalInput uint64

	Outputs        []TransactionAddressInfo
	TotalFCTOutput uint64
	TotalECOutput  uint64

	TxID   string
	Height uint32
	Action [3]bool // Sent, recieved, converted
	Date   string
	Time   string

	//ITrans interfaces.ITransaction
}

// for sorting
type DisplayTransactions []DisplayTransaction

func (slice DisplayTransactions) Len() int {
	return len(slice)
}

func (slice DisplayTransactions) Less(i, j int) bool {
	// Reverse, as higher height = newer
	return slice[i].Height > slice[j].Height
}

func (slice DisplayTransactions) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (w *WalletDB) NewDisplayTransaction(t interfaces.ITransaction) (*DisplayTransaction, error) {
	if t == nil {
		return nil, fmt.Errorf("Transaction is nil")
	}

	_, err := w.TransactionDB.GetAllTXs()
	if err != nil {
		return nil, err
	}

	dt := new(DisplayTransaction)
	//dt.ITrans = t
	dt.TotalInput = 0
	dt.TotalFCTOutput = 0
	dt.TotalECOutput = 0
	dt.Height = t.GetBlockHeight()
	dt.TxID = t.GetHash().String()
	dt.Inputs = make([]TransactionAddressInfo, 0)
	dt.Outputs = make([]TransactionAddressInfo, 0)
	dt.Action = [3]bool{false, false, false}
	dt.Date = t.GetTimestamp().GetTime().Format(("01/02/2006"))
	dt.Time = t.GetTimestamp().GetTime().Format(("15:04:05"))

	ins := t.GetInputs()
	// Inputs
	for _, in := range ins {
		add := primitives.ConvertFctAddressToUserStr(in.GetAddress())
		anp, _ := w.GetGUIAddress(add)
		name := ""
		if anp != nil {
			name = anp.Name
			dt.Action[0] = true
		}

		amt := in.GetAmount()
		dt.TotalInput += amt

		dt.Inputs = append(dt.Inputs, *NewTransactionAddressInfo(name, add, amt, "FCT"))
	}

	outs := t.GetOutputs()
	// FCT Outputs
	for _, out := range outs {
		add := primitives.ConvertFctAddressToUserStr(out.GetAddress())
		anp, _ := w.GetGUIAddress(add)
		name := ""
		if anp != nil {
			name = anp.Name
			dt.Action[1] = true
		}

		amt := out.GetAmount()
		dt.TotalFCTOutput += amt

		dt.Outputs = append(dt.Outputs, *NewTransactionAddressInfo(name, add, amt, "FCT"))
	}

	ecOuts := t.GetECOutputs()
	// EC Outputs
	for _, ecOut := range ecOuts {
		add := primitives.ConvertECAddressToUserStr(ecOut.GetAddress())
		anp, _ := w.GetGUIAddress(add)
		name := ""
		if anp != nil {
			name = anp.Name
			dt.Action[2] = true
		}

		amt := ecOut.GetAmount()
		dt.TotalECOutput += amt

		dt.Outputs = append(dt.Outputs, *NewTransactionAddressInfo(name, add, amt, "EC"))
	}

	return dt, nil
}

// Currently no caching
func (w *WalletDB) GetRelatedTransactions() ([]DisplayTransaction, error) {

	/*block, err := w.Wallet.TXDB().DBO.FetchFBlockHead()
	fmt.Println(block, err)
	if err != nil {
		fmt.Println("Exit 1")
		return err
	}
	// Last update
	start := w.cachedHeight
	heights, err := factom.GetHeights()
	if err != nil {
		fmt.Println("Exit 2")
		return err
	}

	end := heights.LeaderHeight
	_ = end
	_ = start*/
	// TODO: Caching

	// ## No cache solution ##
	// Need to prevent Duplicates
	transMap := make(map[string]interfaces.ITransaction)
	var transList []DisplayTransaction
	adds := w.GetAllGUIAddresses()
	for _, a := range adds {
		transactions, err := w.Wallet.TXDB().GetTXAddress(a.Address)
		if err != nil {
			return nil, err
		}

		for _, trans := range transactions {
			i, _ := transMap[trans.GetHash().String()]
			if i == nil {
				transMap[trans.GetHash().String()] = trans
				dt, err := w.NewDisplayTransaction(trans)
				if err != nil {
					return nil, err
				}
				transList = append(transList, *dt)
			}
		}
	}

	// ## End no cache ##

	// Update to current
	//h := block.GetDBHeight()
	//w.cachedHeight = h

	sort.Sort(DisplayTransactions(transList))
	return transList, nil
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
	err := w.GUIlDB.Put([]byte("gui-wallet"), []byte("wallet"), w.guiWallet)
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
