package raft

import (
	"encoding/binary"
	"errors"
	"github.com/boltdb/bolt"
	"log"
)

var (
	bucketName = []byte("B1")

	KeyNotFoundError = errors.New("not found")
)

type Store struct {
	conn *bolt.DB
	path string
}

func NewStore(path string) () {

}

func (b *Store) initialize() error {
	return
}
func (b *Store) Close() error {

}

func (b *Store) Set(k, v []byte) error {
	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	bucket := tx.Bucket(bucketName)
	if err := bucket.Put(k, v); err != nil {
		return err
	}
	return tx.Commit()
}

func (b *Store) Get(k []byte) ([]byte, error) {
	tx, err := b.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	bucket := tx.Bucket(bucketName)
	v := bucket.Get(k)
	if v == nil {
		return nil, KeyNotFoundError
	}

	return append([]byte(nil), v...), nil
}

func (b *Store) SetUint64(k []byte, v uint64) error {
	return b.Set(k, uint64ToByte(v))
}

func (b *Store) GetUint64(k []byte) (uint64, error) {
	v, err := b.Get(k)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(v), nil
}

func (b *Store) FirstIndex() (uint64, error) {

}

func (b *Store) LastIndex() (uint64, error) {

}

func (b *Store) GetLog(idx uint64, log *log.Logger) error {

}

func (b *Store) StoreLogs(logs []*log.Logger) error {

}

func (b *Store) DeleteRange(min, max uint64) error {

}

func (b *Store) Sync() error {

}

func uint64ToByte(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
