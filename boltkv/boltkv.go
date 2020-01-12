package boltkv

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type BoltItem struct {
	Db *bolt.DB
}

func Open(name string) (*bolt.DB, error) {
	logrus.Infof("[BOLTKV] open db: %s", name)
	db, err := bolt.Open(name, 0666, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (bki *BoltItem) Close() {
	bki.Db.Close()
}

func (bki *BoltItem) Get(k []byte) (value []byte, err error) {
	logrus.Infof("[BOLTKV] get: %s ", k)
	err = bki.Db.View(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket([]byte("bucket")); bucket != nil {
			value = bucket.Get(k)
			return nil
		}
		return errors.New("bucket not found")
	})
	return value, err
}

func (bki *BoltItem) Set(k []byte, v []byte) error {
	logrus.Infof("[BOLTKV] set: %s,%s ", k, v)
	if err := bki.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("bucket"))
		if err != nil {
			return err
		}
		return bucket.Put(k, v)

	}); err != nil {
		logrus.Warnf("[BOLTKV] set failed: %s", err)
		return err
	}
	return nil
}

func (bki *BoltItem) Delete(k []byte) error {
	logrus.Infof("[BOLTKV] delete: %s ", k)
	if err := bki.Db.Update(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket([]byte("bucket")); bucket != nil {
			return bucket.Delete(k)
		}
		return errors.New("bucket not found")
	}); err != nil {
		return err
	}
	return nil
}
