package database

/********************************
 *                              *
 *  Controls LevelDB functions  *
 *                              *
 ********************************/

import (
	"os"
	"sync"

	"github.com/FactomProject/goleveldb/leveldb"
	"github.com/FactomProject/goleveldb/leveldb/opt"
)

type LevelDB struct {
	// lock preventing multiple entry
	dbLock sync.RWMutex
	lDB    *leveldb.DB
}

func (db *LevelDB) Delete(key []byte) error {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	return db.lDB.Delete(key, nil)
}

func (db *LevelDB) Close() error {
	return db.lDB.Close()
}

func (db *LevelDB) Put(key []byte, value []byte) error {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	return db.lDB.Put(key, value, nil)
}

func (db *LevelDB) Get(key []byte) ([]byte, error) {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	return db.lDB.Get(key, nil)
}

func NewLevelDB(filename string) (IDatabase, error) {
	db := new(LevelDB)
	var err error

	var tlDB *leveldb.DB

	_, err = os.Stat(filename)
	if err != nil {
		err = os.Mkdir(filename, 0750)
		if err != nil {
			return nil, err
		}
	}

	opts := &opt.Options{
		Compression: opt.NoCompression,
	}

	tlDB, err = leveldb.OpenFile(filename, opts)
	if err != nil {
		return nil, err
	}
	db.lDB = tlDB

	return db, nil
}
