package database

/********************************
 *                              *
 *          All DB Types        *
 *                              *
 ********************************/

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/database/hybridDB"
	"github.com/FactomProject/factomd/database/mapdb"
)

func NewOrOpenLevelDBWallet(ldbpath string) (interfaces.IDatabase, error) {
	// check if the file exists or if it is a directory
	_, err := os.Stat(ldbpath)

	// create the wallet directory if it doesn't already exist
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(ldbpath), 0777); err != nil {
			fmt.Printf("database error %s\n", err)
		}
	}

	db, err := hybridDB.NewLevelMapHybridDB(ldbpath, false)
	if err != nil {
		fmt.Printf("err opening db: %v\n", err)
	}

	if db == nil {
		fmt.Println("Creating new db ...")
		db, err = hybridDB.NewLevelMapHybridDB(ldbpath, true)

		if err != nil {
			return nil, err
		}
	}
	fmt.Println("Database started from: " + ldbpath)
	return db, nil
}

func NewMapDB() (interfaces.IDatabase, error) {
	return new(mapdb.MapDB), nil
}

func NewOrOpenBoltDBWallet(boltPath string) (interfaces.IDatabase, error) {
	// check if the file exists or if it is a directory
	fileInfo, err := os.Stat(boltPath)
	if err == nil {
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("The path %s is a directory.  Please specify a file name.", boltPath)
		}
	}

	// create the wallet directory if it doesn't already exist
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(boltPath), 0777); err != nil {
			fmt.Printf("database error %s\n", err)
		}
	}

	if err != nil && !os.IsNotExist(err) { //some other error, besides the file not existing
		fmt.Printf("database error %s\n", err)
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Could not use wallet file \"%s\"\n%v\n", boltPath, r)
			os.Exit(1)
		}
	}()
	db := hybridDB.NewBoltMapHybridDB(nil, boltPath)

	fmt.Println("Database started from: " + boltPath)
	return db, nil
}
