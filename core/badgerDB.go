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
func (b *BDB) get(key []byte) ([]byte, error) {

	value := []byte{}
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	return value, err
}
func (b *BDB) set(key, value []byte) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
func (b *BDB) delete(key []byte) (bool, error) {
	err := b.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete(key); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
func (b *BDB) snapShotItems() <-chan DataItem {
	s := make(chan DataItem)
	return s
}
