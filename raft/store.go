package raft

import (
	"encoding/binary"
	"errors"
	"github.com/boltdb/bolt"
	"log"
)

const defaultDBFileMode = 0600

var (
	logsName   = []byte("logs")
	bucketName = []byte("bucket")

	KeyNotFoundError = errors.New("not found")
)

type Store struct {
	conn *bolt.DB
	path string
}

type Options struct {
	BoltOptions *bolt.Options
	Path        string
	NoSync      bool
}

func (o *Options) isReadOnly() bool {
	return o != nil && o.BoltOptions != nil && o.BoltOptions.ReadOnly
}

func NewStore(path string) (*Store, error) {
	return New(Options{
		BoltOptions: nil,
		Path:        path,
		NoSync:      false,
	})
}

//New use the open BoltDB and prepare for raft backend
func New(options Options) (*Store, error) {
	DB, err := bolt.Open(options.Path, defaultDBFileMode, options.BoltOptions)
	if err != nil {
		return nil, err
	}
	DB.NoSync = options.NoSync

	store := &Store{
		conn: DB,
		path: options.Path,
	}

	if !options.isReadOnly() {
		if err := store.initialize(); err != nil {
			store.Close()
			return nil, err
		}
	}
	return store, nil
}

func (b *Store) initialize() error {
	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists(logsName); err != nil {
		return err
	}
	if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
		return err
	}

	return tx.Commit()
}
func (b *Store) Close() error {
	return b.conn.Close()
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
