package database

// General database interface

type IDatabase interface {
	Close() error
	Delete(key []byte) error
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
}
