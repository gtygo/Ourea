package core

import (
	"github.com/dgraph-io/badger"
	"log"
	"os"
)

type BDB struct {
	DBdir    string
	valueDir string
	db       *badger.DB
	logger   *log.Logger
}

func NewDB(DBdir, valueDir string) (*BDB, error) {
	opts := badger.DefaultOptions(DBdir)
	opts.ValueDir = valueDir
	opts.SyncWrites = false
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BDB{
		DBdir:    DBdir,
		valueDir: valueDir,
		db:       db,
		logger:   log.New(os.Stderr, "[db] ", log.LstdFlags),
	}, nil
}
func (b *BDB) Get(key []byte) ([]byte, error) {
	return nil, nil
}
func (b *BDB) Set(key, value []byte) error {
	return nil
}
func (b *BDB) Delete(key []byte) (bool, error) {
	return false, nil
}
func (b *BDB) SnapShotItems() <-chan DataItem {
	s := make(chan DataItem)
	return s
}
