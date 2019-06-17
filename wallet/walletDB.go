// .
//
// WalletDB
//
// The wallet consists of three 3 main structures that are interacted with:
// Wallet : From 'factom' library, this contains all the funcationality related with transactions.
// GuiWallet : Contains all the addresses associated with the addresses in the Wallet.
// TransactionDB : Contains all transactions from factom. Has helpful functions to find transactions.
package wallet

/*
 * Manages all the addresses and 2 databases (Wallet DB and GUI DB)
 */

import (
	"fmt"
	"os"
	"os/user"
	"sort"
	"strings"
	"sync"

	"github.com/FactomProject/enterprise-wallet/address"
	"github.com/FactomProject/enterprise-wallet/wallet/database"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factomd/common/primitives"

	// "github.com/FactomProject/factom/wallet/wsapi"
	"encoding/json"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/database/mapdb"
	"github.com/FactomProject/go-bip39"
)

// List of int to db type
const (
	MAP int = iota
	LDB
	BOLT
	ENCRYPTED
)

// Default settings
var (
	GUI_DB    = MAP
	WALLET_DB = MAP
	TX_DB     = MAP
)

// This has to do with launching via CLI. How often to show progress in syncing
var (
	STEPS_TO_PRINT int = 10000 // How many steps needed to alert user of progress
)

// WalletDB interacting with LDB and factom/wallet
//   The LDB doesn't need to be updated often, so we save after every add and only
//   deal with cached version
type WalletDB struct {
	GUIlDB        interfaces.IDatabase      //database.IDatabase        // GUI DB
	guiWallet     *WalletStruct             // Cached version on GUI LDB
	Wallet        *wallet.Wallet            // Wallet from factom/wallet
	TransactionDB *wallet.TXDatabaseOverlay // Used to display transactions

	// Used to cache related transactions
	// This is rebuilt upon every launch
	relatedTransactionLock   sync.RWMutex                       // For all variables associated with related transaction caching
	cachedTransactions       []DisplayTransaction               // All sorted transactions already found
	ActiveCachedTransactions []DisplayTransaction               // Active cache being used.
	cachedHeight             uint32                             // Last FBlock height used
	transMap                 map[string]DisplayTransaction      // Prevent duplicate transactions
	addrMap                  map[string]address.AddressNamePair // Find addresses quick, All addresses already searched for up to last FBlock

	changeAddrMap     map[string]UpdateANP // Addresses that need to be moved to addrmap
	changeAddrMapLock sync.Mutex

	BalanceCache *BalanceCache

	quit bool
}

// LoadWalletDB is the same as New
func LoadWalletDB(v1Import bool, password string) (*WalletDB, error) {
	return NewWalletDB(v1Import, password)
}

func NewWalletDB(v1Import bool, password string) (*WalletDB, error) {
	w := new(WalletDB)

	if WALLET_DB == ENCRYPTED {
		GUI_DB = ENCRYPTED
	}

	var db interfaces.IDatabase
	var err error
	switch GUI_DB { // Decides type of wallet DB
	case MAP:
		db, err = database.NewMapDB()
	case LDB:
		db, err = database.NewOrOpenLevelDBWallet(GetHomeDir() + guiLDBPath)
	case BOLT:
		db, err = database.NewOrOpenBoltDBWallet(GetHomeDir() + guiBoltPath)
	case ENCRYPTED:
		db, err = database.NewOrOpenBoltDBWallet(GetHomeDir() + guiEncryptedBoltPath)
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

	var wal *wallet.Wallet

	// Cannot import to encrypted
	if WALLET_DB == ENCRYPTED {
		v1Import = false
	}

	switch v1Import {
	case true:
		if WALLET_DB == MAP {
			// Let fallthrough
		} else {
			m2Path := ""
			// If true, and there is no M2 file, we will check for M1 file before making a new
			// If a map, then we ignore and do the else
			switch WALLET_DB { // Decides type of wallet DB
			case LDB:
				m2Path = GetHomeDir() + walletLDBPath
				_, err = os.Stat(GetHomeDir() + walletLDBPath)
			case BOLT:
				m2Path = GetHomeDir() + walletBoltPath
				_, err = os.Stat(GetHomeDir() + walletBoltPath)
			case ENCRYPTED:
				return nil, fmt.Errorf("Cannot import from v1 to encrypted wallet")
			}
			if err != nil { // No M2 file, lets grab from M1
				m1Path := ""
				if WalletBoltV1Path == "/.factom/factoid_wallet_bolt.db" {
					m1Path = GetHomeDir() + WalletBoltV1Path
				} else {
					m1Path = WalletBoltV1Path
				}
				_, err = os.Stat(m1Path)
				if err != nil { // No M1 file, lets go as normal
					// Let fallthrough
				} else { // M1 file found, no M2 file. Let's import
					fmt.Printf("Importing from M1 Wallet at %s....\n", m1Path)
					wal, err = wallet.ImportV1Wallet(m1Path, m2Path)
					break // We got a wal file, let's break
				}
			} else { // There is an M2 file, we go as normal
				// Let falltrhough
			}
		}
		fallthrough
	case false:
		switch WALLET_DB {
		case MAP:
			wal, err = wallet.NewMapDBWallet()
		case LDB:
			wal, err = wallet.NewOrOpenLevelDBWallet(GetHomeDir() + walletLDBPath)
		case BOLT:
			wal, err = wallet.NewOrOpenBoltDBWallet(GetHomeDir() + walletBoltPath)
		case ENCRYPTED:
			wal, err = wallet.NewEncryptedBoltDBWallet(GetHomeDir()+walletEncryptedBoltPath, password)
		}
	}
	if err != nil {
		return nil, err
	}

	w.Wallet = wal

	var txdb *wallet.TXDatabaseOverlay
	switch TX_DB {
	case MAP:
		txdb = wallet.NewTXOverlay(new(mapdb.MapDB))
		err = nil
	case LDB:
		txdb, err = wallet.NewTXLevelDB(GetHomeDir() + txdbLDBPath)
	case BOLT:
		txdb, err = wallet.NewTXBoltDB(GetHomeDir() + txdbBoltPath)
	}

	if err != nil {
		return nil, fmt.Errorf("Could not add transaction database to wallet: %s\n", err.Error())
	}

	w.Wallet.AddTXDB(txdb)

	w.TransactionDB = w.Wallet.TXDB()
	if w.TransactionDB != nil { // Update DB
		//w.TransactionDB.GetAllTXs()
	}

	err = w.UpdateGUIDB()
	if err != nil {
		return nil, err
	}

	w.transMap = make(map[string]DisplayTransaction)
	w.addrMap = make(map[string]address.AddressNamePair)
	w.changeAddrMap = make(map[string]UpdateANP)
	w.cachedHeight = 0
	w.ActiveCachedTransactions = w.cachedTransactions

	w.BalanceCache = NewBalanceCache()

	return w, nil
}

// DisplayTransactions is used for sorting
type DisplayTransactions []DisplayTransaction

func (slice DisplayTransactions) Len() int {
	return len(slice)
}

func (slice DisplayTransactions) Less(i, j int) bool {
	return !slice[i].ExactTime.Before(slice[j].ExactTime)
}

func (slice DisplayTransactions) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice DisplayTransactions) IsSameAs(comp DisplayTransactions) bool {
	for i := 0; i < slice.Len(); i++ {
		if !slice[i].IsSameAs(comp[i]) {
			return false
		}
	}
	return true
}

func (slice DisplayTransactions) IsSimilarTo(comp DisplayTransactions) bool {
	for i := 0; i < slice.Len(); i++ {
		if !slice[i].IsSimilarTo(comp[i]) {
			return false
		}
	}
	return true
}

// NewDisplayTransaction is used because we don't keep the original interface of a transaction, but
// build our own
func (w *WalletDB) NewDisplayTransaction(t interfaces.ITransaction) (*DisplayTransaction, error) {
	if t == nil {
		return nil, fmt.Errorf("Transaction is nil")
	}

	dt := new(DisplayTransaction)
	//dt.ITrans = t
	dt.TotalInput = 0
	dt.TotalFCTOutput = 0
	dt.TotalECOutput = 0
	dt.Height = t.GetBlockHeight()
	dt.TxID = t.GetSigHash().String()
	dt.Inputs = make([]TransactionAddressInfo, 0)
	dt.Outputs = make([]TransactionAddressInfo, 0)
	dt.Action = [3]bool{false, false, false}
	dt.ExactTime = t.GetTimestamp().GetTime()
	dt.Date = dt.ExactTime.Format(("01/02/2006"))
	dt.Time = dt.ExactTime.Format(("15:04:05"))
	ins := t.GetInputs()
	// Inputs
	for _, in := range ins {
		add := primitives.ConvertFctAddressToUserStr(in.GetAddress())
		//anp, _ := w.GetGUIAddress(add)
		anp, ok := w.addrMap[add]
		name := ""
		if ok {
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
		//anp, _ := w.GetGUIAddress(add)
		anp, ok := w.addrMap[add]
		name := ""
		if ok {
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
		//anp, _ := w.GetGUIAddress(add)
		anp, ok := w.addrMap[add]
		name := ""
		if ok {
			name = anp.Name
			dt.Action[2] = true
		}

		amt := ecOut.GetAmount()
		dt.TotalECOutput += amt

		dt.Outputs = append(dt.Outputs, *NewTransactionAddressInfo(name, add, amt, "EC"))
	}
	return dt, nil
}

func (w *WalletDB) ExportSeed() (string, error) {
	return w.Wallet.GetSeed()
}

var PROCESSING_RELATED_TRANSACTIONS = false

func prtOff() {
	PROCESSING_RELATED_TRANSACTIONS = false
}

// The user's need feedback on the sync. The function prints out the best information,
// but to display to front end, we need a simple 1 percent number
// 	Stage 0: Setup
//	Stage 1: Gathering Transactions
//	Stage 2: Checking New Addresses
//	Stage 3: Sorting
var RTLock sync.RWMutex
var RTStage int = 0

func (w *WalletDB) SetStage(stage int) {
	RTLock.Lock()
	RTStage = stage
	RTLock.Unlock()
}

func (w *WalletDB) GetStage() int {
	RTLock.RLock()
	t := RTStage
	RTLock.RUnlock()
	return t
}

// GetRelatedTransactions
// This function grabs all transactions related to any address in the address book
// and sorts them by time.Time. If a new address is added, this will grab all transactions
// from that new address and insert them.
func (w *WalletDB) GetRelatedTransactions() (dt []DisplayTransaction, err error) {
	if PROCESSING_RELATED_TRANSACTIONS { // Already working on it
		return
	}

	w.SetStage(0)
	if w.quit {
		return nil, fmt.Errorf("Wallet is closing, stoped fetching transactions")
	}

	// If we print 1 step, we should print all so user knows it is done
	// Some steps may be very quick
	printSteps := false

	PROCESSING_RELATED_TRANSACTIONS = true
	defer prtOff()

	// Temporary
	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if recover() != nil {
			err = fmt.Errorf("There was an issue trying to load the database. Please try again in a few seconds. If you keep encountering this error," +
				"factomd might be having issues syncing with the network.")
		}
	}()
	w.relatedTransactionLock.Lock()
	defer w.relatedTransactionLock.Unlock()

	// Get current Fblock height
	var i int
	var block interfaces.IFBlock
	for i = 0; i < 2; i++ { // 2 tries, if fails first, updates transactions and trys again
		block, err = w.TransactionDB.DBO.FetchFBlockHead()
		if err != nil {
			return nil, fmt.Errorf("There has been an issue when trying to grab a block from the database. Please try again in a few seconds.")
		}
		if block == nil {
			if i == 0 {
				w.TransactionDB.GetAllTXs()
			} else {
				return nil, fmt.Errorf("Error with loading transaction database. It could be in the process of loading all transactions. Try waiting a minute and reloading the page.")
			}
		} else {
			break
		}
	}

	if block.GetDatabaseHeight() == 0 {
		// return nil, fmt.Errorf("Must wait 1 block and try again.")
	}

	var oldHeight uint32
	if block != nil {
		oldHeight = w.cachedHeight
		w.cachedHeight = block.GetDatabaseHeight()
	} else {
		w.TransactionDB.GetAllTXs() // UpdateDB for next attempt if user tries again
		return nil, fmt.Errorf("Error with loading transaction database. Try waiting a minute and reloading the page.")
	}

	if w.quit {
		return nil, fmt.Errorf("Wallet is closing, stoped fetching transactions")
	}

	//
	// STAGE 1
	w.SetStage(1)
	// Gather all new transactions, and add the ones related to us into a list
	//

	// Get all new transaction to go through
	transactions, err := w.TransactionDB.GetTXRange(int(oldHeight), int(w.cachedHeight))
	if err != nil {
		return nil, fmt.Errorf("Error with loading transactions from the database. When grabbing from transaction range.")
	}
	totalTransactions := len(transactions)
	var newTransactions []DisplayTransaction
	// Sort through new transactions for any related
	for i, trans := range transactions {
		if totalTransactions > STEPS_TO_PRINT && i%STEPS_TO_PRINT == 0 {
			fmt.Printf("Step 1/3 for Transactions %d / %d\n", i, totalTransactions)
		}
		added := false
		for i = 0; i < 3; i++ {
			var addresses []string
			switch i {
			case 0:
				addrs := trans.GetInputs()
				for _, a := range addrs {
					addresses = append(addresses, primitives.ConvertFctAddressToUserStr(a.GetAddress()))
				}
			case 1:
				addrs := trans.GetOutputs()
				for _, a := range addrs {
					addresses = append(addresses, primitives.ConvertFctAddressToUserStr(a.GetAddress()))
				}
			case 2:
				addrs := trans.GetECOutputs()
				for _, a := range addrs {
					addresses = append(addresses, primitives.ConvertECAddressToUserStr(a.GetAddress()))
				}
			}

			for _, addr := range addresses { // If it makes through this loop will check next set of addresses
				_, ok := w.addrMap[addr]
				if ok {
					dt, err := w.NewDisplayTransaction(trans)
					if err != nil {
						break // Error with transaction
					}

					_, ok := w.transMap[dt.TxID]
					if !ok {
						newTransactions = append(newTransactions, *dt)
						w.transMap[dt.TxID] = *dt
					}
					added = true
					break // Transaction added
				}
			}

			if added {
				break // Transaction added, break out of this transaction
			}
		}
	}

	if totalTransactions > STEPS_TO_PRINT || printSteps {
		printSteps = true
		fmt.Printf("Step 1/3 for Transactions %d / %d\n", totalTransactions, totalTransactions)
	}

	if w.quit {
		return nil, fmt.Errorf("Wallet is closing, stoped fetching transactions")
	}

	//
	// STAGE 2
	w.SetStage(2)
	// Sort transactions we already have, and get transactions from any new addresses
	//

	// Sort the new ones
	sort.Sort(DisplayTransactions(newTransactions))

	// Prepend them to the old cache
	w.cachedTransactions = append(newTransactions, w.cachedTransactions...)
	// Find all new addresses, need to do additional handling and inserting
	var moreTransactions []DisplayTransaction
	anps := w.GetAllMyGUIAddresses()
	var newAddrs []string
	totalTransactions = 0
	currentCheckpoint := 0
	for _, a := range anps {
		if w.quit {
			return nil, fmt.Errorf("Wallet is closing, stoped fetching transactions")
		}

		_, ok := w.addrMap[a.Address]
		if ok { // Found

		} else { // New addr
			w.addrMap[a.Address] = a
			newAddrs = append(newAddrs, a.Address)
			trans, err := w.TransactionDB.GetTXAddress(a.Address)
			if err == nil {
				if len(trans) > 0 {
					totalTransactions += len(trans)
					// This takes some real time for huge amounts
					for _, t := range trans {
						currentCheckpoint++
						if totalTransactions > STEPS_TO_PRINT && currentCheckpoint%STEPS_TO_PRINT == 0 {
							fmt.Printf("Step 2/3 for Transactions %d / %d\n", i+currentCheckpoint, totalTransactions)
						}
						dt, _ := w.NewDisplayTransaction(t)
						moreTransactions = append(moreTransactions, *dt)
					}
					//moreTransactions = append(moreTransactions, trans...)
				}
			}
			currentCheckpoint = totalTransactions
		}
	}
	if totalTransactions > 1000 || printSteps {
		printSteps = true
		fmt.Printf("Step 2/3 for Transactions %d / %d\n", totalTransactions, totalTransactions)
	}

	if w.quit {
		return nil, fmt.Errorf("Wallet is closing, stoped fetching transactions")
	}

	//
	// STAGE 3
	w.SetStage(3)
	// Insert any new transactions from new addresses into our list to append
	//

	totalTransactions = len(moreTransactions)
	/* This to end of function breaks the attempt to build for windows for some reason */
	// Binary search and insert new transactions from new addresses
	for i, t := range moreTransactions {
		if w.quit {
			return nil, fmt.Errorf("Wallet is closing, stoped fetching transactions")
		}

		if totalTransactions > STEPS_TO_PRINT && i%STEPS_TO_PRINT == 0 {
			fmt.Printf("Step 3/3 for Transactions %d / %d\n", i, totalTransactions)
		}
		if _, ok := w.transMap[t.TxID]; ok {
			continue
		}

		i = w.findTransactionIndex(t)

		if i < len(w.cachedTransactions) && w.cachedTransactions[i].TxID == t.TxID {
			// t is present at w.cachedTransactions[i], already there. We need to update the 'Actions'
			// field. If we have the same transaction as before, we don't need to add a new one, but update
			// the existing
			for counter := 0; counter < 3; counter++ {
				// If one or other is true, we want to keep that
				w.cachedTransactions[i].Action[counter] = w.cachedTransactions[i].Action[counter] || t.Action[counter]
			}
		} else {
			// t is not present in w.cachedTransactions,
			// but i is the index where it would be inserted.
			w.transMap[t.TxID] = t // Add to cache
			// Insert
			w.cachedTransactions = append(w.cachedTransactions[:i], append([]DisplayTransaction{t}, w.cachedTransactions[i:]...)...)
		}
	}
	if totalTransactions > STEPS_TO_PRINT || printSteps {
		printSteps = true
		fmt.Printf("Step 3/3 for Transactions %d / %d\n", totalTransactions, totalTransactions)
		fmt.Printf("Finishing up sync....\n")
	}

	// The edge case of no transactions. If you have no related transactions, we still need to signal we
	// are completely loaded. So we will add a blank transaction with an "empty" txid, which is impossible to get otherwise.
	if len(w.cachedTransactions) == 0 {
		empty := new(DisplayTransaction)
		empty.TxID = "empty"
		var temp []DisplayTransaction
		temp = append(temp, *empty)
		return temp, nil
	} else {
		return w.cachedTransactions, nil
	}
}

// findTransactionIndex uses binary search
func (w *WalletDB) findTransactionIndex(t DisplayTransaction) int {
	low := 0
	high := len(w.cachedTransactions) - 1

	for low <= high {
		mid := low + ((high - low) / 2)
		if w.cachedTransactions[mid].TxID == t.TxID {
			return mid
		}
		if !w.cachedTransactions[mid].ExactTime.Before(t.ExactTime) {
			//high = mid - 1
			low = mid + 1
		} else {
			//low = mid + 1
			high = mid - 1
		}
	}

	return low
}

// GetRelatedTransactionsNoCaching is the no cache solution, not going to use it. It is too slow, but was used in early phases and kept
// for testing comparisons as this should be all inclusive and correct
func (w *WalletDB) GetRelatedTransactionsNoCaching() ([]DisplayTransaction, error) {
	// ## No cache solution ##
	transMap := make(map[string]interfaces.ITransaction)
	var transList []DisplayTransaction
	adds := w.GetAllMyGUIAddresses()
	for _, a := range adds {
		transactions, err := w.TransactionDB.GetTXAddress(a.Address)
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

	sort.Sort(DisplayTransactions(transList))
	return transList, nil
}

func (w *WalletDB) GetGUIWalletJSON(getBals bool) ([]byte, error) {
	if getBals {
		w.AddBalancesToAddresses()
	}
	return json.Marshal(w.guiWallet)
}

func (w *WalletDB) AddBalancesToAddresses() {
	w.guiWallet.AddBalancesToAddresses()
}

// UpdateGUIDB grabs the list of addresses from the walletDB and updates our
// GUI with any that are missing. All will be external
func (w *WalletDB) UpdateGUIDB() error {
	faAdds, ecAdds, err := w.Wallet.GetAllAddresses()
	if err != nil {
		return err
	}

	var addMap map[string]string
	addMap = make(map[string]string)

	var names []string
	var addresses []string

	guiAdds := w.GetAllMyGUIAddresses()

	// Add addresses to GUI from cli
	for _, fa := range faAdds {
		_, list := w.GetGUIAddress(fa.String())
		if list == -1 {
			names = append(names, "FA-Imported-From-CLI")
			addresses = append(addresses, fa.String())
		}
		addMap[fa.String()] = fa.String()
	}

	for _, ec := range ecAdds {
		_, list := w.GetGUIAddress(ec.String())
		if list == -1 {
			names = append(names, "EC-Imported-From-CLI")
			addresses = append(addresses, ec.String())
		}
		addMap[ec.String()] = ec.String()
	}

	// Add in new guys
	if len(names) > 0 {
		err = w.addBatchGUIAddresses(names, addresses)
		if err != nil {
			return err
		}
	}

	// Missing from CLI? We need to remove them here
	for _, guiAdd := range guiAdds {
		if _, ok := addMap[guiAdd.Address]; !ok {
			w.RemoveAddressFromAnyList(guiAdd.Address)
		}
	}

	return w.Save()
}

func (w *WalletDB) Close() error {
	w.quit = true

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

	err = w.TransactionDB.Close()
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
		return fmt.Errorf("The GUI database encountered an error while saving. Your last action may not have been saved. Shutting down the GUI may lose your last action.")
		//return err
	}

	return nil
}

func (w *WalletDB) GenerateFactoidAddress(name string) (*address.AddressNamePair, error) {
	address, err := w.Wallet.GenerateFCTAddress()

	if err != nil {
		return nil, fmt.Errorf("There has been an error generating a new Factoid address, please try again.")
		//return nil, err
	}

	anp, err := w.guiWallet.AddSeededAddress(name, address.String(), 1)
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

	anp, err := w.guiWallet.AddSeededAddress(name, address.String(), 2)
	if err != nil {
		return nil, err
	}

	w.Save()
	if err != nil {
		return nil, err
	}

	return anp, nil
}

func (w *WalletDB) RemoveAddress(address string, list int) (*address.AddressNamePair, error) {
	anp, _, _ := w.guiWallet.GetAddress(address)

	_, err := w.guiWallet.RemoveAddress(anp.Address, list)
	if err != nil {
		return nil, err
	}

	err = w.Save()
	if err != nil {
		return nil, err
	}

	return anp, nil
}

func (w *WalletDB) RemoveAddressFromAnyList(address string) (*address.AddressNamePair, error) {
	anp, err := w.guiWallet.RemoveAddressFromAnyList(address)
	if err != nil {
		return nil, err
	}

	err = w.Save()
	if err != nil {
		return nil, err
	}

	return anp, nil
}

func (w *WalletDB) AddExternalAddress(name string, public string) (*address.AddressNamePair, error) {
	if !factom.IsValidAddress(public) {
		return nil, fmt.Errorf("Not a valid public key")
	}

	anp, err := w.addGUIAddress(name, public, 3)
	if err != nil {
		return nil, err
	}

	return anp, nil
}

func (w *WalletDB) ImportSeed(seed string) error {
	seedStruct := new(wallet.DBSeed)
	seedStruct.MnemonicSeed = seed

	v := bip39.IsMnemonicValid(seed)
	if !v {
		return fmt.Errorf("not a valid seed")
	}

	err := w.Wallet.InsertDBSeed(seedStruct)
	if err != nil {
		return err
	}

	w.guiWallet.ResetSeeded()
	w.UpdateGUIDB()
	return nil
}

func (w *WalletDB) ImportKoinify(name string, koinify string) (*address.AddressNamePair, error) {
	add, err := factom.MakeFactoidAddressFromKoinify(koinify)
	if err != nil {
		return nil, err
	}

	err = w.Wallet.InsertFCTAddress(add)
	if err != nil {
		return nil, err
	}

	anp, err := w.addGUIAddress(name, add.String(), 1)
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
			return nil, fmt.Errorf("There has been an error when trying to convert the secret key to an address, the secret key may not be valid.")
			//return nil, err
		}

		err = w.Wallet.InsertFCTAddress(add)
		if err != nil {
			return nil, fmt.Errorf("There has been an error trying to insert the new address into the wallet. Please try again.")
			//return nil, err
		}

		anp, err := w.addGUIAddress(name, add.String(), 1)
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

		anp, err := w.addGUIAddress(name, add.String(), 2)
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

// addBatchGUIAddresses only adds to GUI Database
func (w *WalletDB) addBatchGUIAddresses(names []string, addresses []string) error {
	if len(names) != len(addresses) {
		return fmt.Errorf("List length does not match")
	}

	for i := 0; i < len(names); i++ {
		if addresses[i][:2] == "FA" {
			w.addGUIAddress(names[i], addresses[i], 1)
		} else {
			w.addGUIAddress(names[i], addresses[i], 2)
		}
	}

	return w.Save()
}

// addGUIAddress only adds to GUI database
func (w *WalletDB) addGUIAddress(name string, addressStr string, list int) (*address.AddressNamePair, error) {
	var anp *address.AddressNamePair
	var err error
	if list <= 0 || list > 3 {
		return nil, fmt.Errorf("Invalid list")
	}
	if addressStr[:2] == "FA" {
		if list == 2 {
			return nil, fmt.Errorf("Factoid address cannot go in Entry credit list")
		}
		anp, err = w.guiWallet.AddAddress(name, addressStr, list)
	} else {
		if list == 1 {
			return nil, fmt.Errorf("Entry credit address cannot go in Factoid list")
		}
		anp, err = w.guiWallet.AddAddress(name, addressStr, list)
	}

	if err != nil {
		return nil, err
	}
	err = w.Save()
	if err != nil {
		return nil, err
	}

	return anp, nil
}

// GetGUIAddress returns address with associated name
// List is 0 for not found, 1 for Factoid address, 2 for EC Address, 3 for External
func (w *WalletDB) GetGUIAddress(address string) (anp *address.AddressNamePair, list int) {
	anp, list, _ = w.guiWallet.GetAddress(address)
	return
}

// ScrubDisplayTransactionsForNameChanges scrubs all transactions before serving to front end. Changes the names to
// the current names of the addresses, as user can change the name of their addresses.
func (w *WalletDB) ScrubDisplayTransactionsForNameChanges(list []DisplayTransaction) []DisplayTransaction {
	w.relatedTransactionLock.Lock()
	w.getNameChanges()

	for i := range list {
		ins := list[i].Inputs
		for ii := range ins {
			add := list[i].Inputs[ii].Address
			anp, ok := w.addrMap[add]
			if ok {
				list[i].Inputs[ii].Name = anp.Name
				list[i].Action[0] = true
			}
		}
		outs := list[i].Outputs
		for ii := range outs {
			add := list[i].Outputs[ii].Address
			anp, ok := w.addrMap[add]
			if ok {
				list[i].Outputs[ii].Name = anp.Name
				if add[:2] == "FA" {
					list[i].Action[1] = true
				} else {
					list[i].Action[2] = true
				}
			}
		}
	}
	w.relatedTransactionLock.Unlock()

	return list
}

type UpdateANP struct {
	Address string
	Name    string
}

func (w *WalletDB) getNameChanges() {
	w.changeAddrMapLock.Lock()
	for k, v := range w.changeAddrMap {
		anp := w.addrMap[k]
		anp.Name = v.Name
		w.addrMap[k] = anp
	}
	w.changeAddrMapLock.Unlock()
}

func (w *WalletDB) ChangeAddressName(address string, toName string) error {
	err := w.guiWallet.ChangeAddressName(address, toName)
	if err != nil {
		return err
	}

	w.changeAddrMapLock.Lock() // Related Transactions uses this
	w.changeAddrMap[address] = UpdateANP{Address: address, Name: toName}
	w.changeAddrMapLock.Unlock()
	return w.Save()
}

func (w *WalletDB) GetTotalGUIAddresses() uint64 {
	return w.guiWallet.GetTotalAddressCount()
}

func (w *WalletDB) GetAllGUIAddresses() []address.AddressNamePair {
	return w.guiWallet.GetAllAddresses()
}

func (w *WalletDB) GetAllMyGUIAddresses() []address.AddressNamePair {
	return w.guiWallet.GetAllMyGUIAddresses()
}

func (w *WalletDB) IsValidAddress(address string) bool {
	return factom.IsValidAddress(address)
}

func (w *WalletDB) GetECBalance() int64 {
	w.guiWallet.RLock()
	defer w.guiWallet.RUnlock()
	return w.guiWallet.ECTotal
}

func (w *WalletDB) GetFactoidBalance() int64 {
	w.guiWallet.RLock()
	defer w.guiWallet.RUnlock()
	return w.guiWallet.FactoidTotal
}

func (w *WalletDB) FactomdOnline() (bool, string) {
	_, err := factom.GetHeights()
	if err != nil {
		return false, factom.FactomdServer()
	} else {
		return true, factom.FactomdServer()
	}
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
