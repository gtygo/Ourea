package boltkv

import (
	"github.com/boltdb/bolt"
)

type BoltItem struct {
	Db *bolt.DB
}

func Open(name string) (*bolt.DB, error) {
	db, err := bolt.Open(name, 0666, nil)
	if err != nil {
		return nil, err
	}

	println("open")
	return db, nil
}

func (bki *BoltItem) Close() {
	bki.Db.Close()
}

func (bki *BoltItem) Get(k []byte) ([]byte, error) {
	println("get")
	return nil, nil
}

func (bki *BoltItem) Set(k []byte, v []byte) error {
	println("set")
	return nil
}

func (bki *BoltItem) Delete(k []byte) error {
	println("delete")
	return nil
}
