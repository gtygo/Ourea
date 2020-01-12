package boltkv

import (
	"log"

	"github.com/boltdb/bolt"
)

type BoltItem struct {
	Db *bolt.DB
}

func Open(name string) (*bolt.DB, error) {
	log.Printf("[BOLTKV] open db: %s", name)
	db, err := bolt.Open(name, 0666, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (bki *BoltItem) Close() {
	bki.Db.Close()
}

func (bki *BoltItem) Get(k []byte) ([]byte, error) {
	log.Printf("[BOLTKV] get: %s ", k)
	return nil, nil
}

func (bki *BoltItem) Set(k []byte, v []byte) error {
	log.Printf("[BOLTKV] set: %s,%s ", k, v)
	return nil
}

func (bki *BoltItem) Delete(k []byte) error {
	log.Printf("[BOLTKV] delete: %s ", k)
	return nil
}
