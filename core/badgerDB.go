package core

import (
	"errors"
	"github.com/dgraph-io/badger"
	"log"
	"os"
)

var ErrIterFinished = errors.New("ERR iteration finished successfully")

type BDB struct {
	DBdir    string
	valueDir string
	db       *badger.DB
	logger   *log.Logger
}

type KVItem struct {
	key []byte
	val []byte
	err error
}

func (i *KVItem) IsFinished() bool {
	return i.err == ErrIterFinished
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
	// create a new channel
	ch := make(chan DataItem, 1024)

	// generate items from snapshot to channel
	go b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		keyCount := 0
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.ValueCopy(nil)

			kvi := &KVItem{
				key: append([]byte{}, k...),
				val: append([]byte{}, v...),
				err: err,
			}

			// write kvitem to channel with last error
			ch <- kvi
			keyCount = keyCount + 1

			if err != nil {
				return err
			}
		}

		// just use nil kvitem to mark the end
		kvi := &KVItem{
			key: nil,
			val: nil,
			err: ErrIterFinished,
		}
		ch <- kvi

		b.logger.Printf("Snapshot total %d keys", keyCount)

		return nil
	})

	// return channel to persist
	return ch
}

func (b *BDB) close() {
	b.db.Close()
}
