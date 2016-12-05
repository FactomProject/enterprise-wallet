package database

/********************************
 *                              *
 *    Controls Map functions    *
 *                              *
 ********************************/
/*
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type MapDB struct {
	// lock preventing multiple entry
	dbLock sync.RWMutex
	mapDB  map[string][]byte
}

func getKeyString(key []byte) string {
	hash := sha256.Sum256(key)
	keyStr := hex.EncodeToString(hash[:])

	return keyStr
}

func (db *MapDB) Delete(key []byte) error {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	delete(db.mapDB, getKeyString(key))

	return nil
}

func (db *MapDB) Close() error {
	return nil
}

func (db *MapDB) Put(key []byte, value []byte) error {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	db.mapDB[getKeyString(key)] = value

	return nil
}

func (db *MapDB) Get(key []byte) ([]byte, error) {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	val, _ := db.mapDB[getKeyString(key)]

	if val == nil {
		return nil, fmt.Errorf("No value exists")
	}

	return val, nil
}

func NewMapDB() (IDatabase, error) {
	db := new(MapDB)
	db.mapDB = make(map[string][]byte)

	return db, nil
}*/
