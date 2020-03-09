package boltkv

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type BoltItem struct {
	Db *bolt.DB
}

var ErrBucketNotFound =errors.New("bucket not found")

func Open(name string) (*BoltItem, error) {
	logrus.Infof("[BOLTKV] open db: %s", name)
	db, err := bolt.Open(name, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &BoltItem{Db:db}, nil
}

func (bki *BoltItem) Close() {
	bki.Db.Close()
}

func (bki *BoltItem) Get() (m map[string]interface{}, err error) {
	m=make(map[string]interface{})
	err = bki.Db.View(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket([]byte("bucket")); bucket != nil {

			if err := bucket.ForEach(func(k, v []byte) error {
				m[string(k)]=string(v)
				return nil
			}); err != nil {
				return err
			}
			return nil
		}
		return ErrBucketNotFound
	})
	return m, err
}

func (bki *BoltItem) Set(items map[string]interface{}) error {
	if err := bki.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("bucket"))
		if err != nil {
			return err
		}

		for k,_:=range items{
			v:=items[k].(string)
			if err:=bucket.Put([]byte(k),[]byte(v));err!=nil{
				return err
			}
		}
		return nil
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
