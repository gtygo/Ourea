package raft

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/ugorji/go/codec"
)

const defaultDBFileMode = 0600

var (
	logsName   = []byte("logs")
	bucketName = []byte("bucket")

	KeyNotFoundError = errors.New("not found")
)

type Store struct {
	db   *bolt.DB
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
		db:   DB,
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
	tx, err := b.db.Begin(true)
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
	return b.db.Close()
}

func (b *Store) Set(k, v []byte) error {
	tx, err := b.db.Begin(true)
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
	tx, err := b.db.Begin(false)
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
	tx, err := b.db.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	cursor := tx.Bucket(logsName).Cursor()
	if first, _ := cursor.First(); first == nil {
		return 0, nil
	} else {
		return bytesToUint64(first), nil
	}
}

func (b *Store) LastIndex() (uint64, error) {
	tx, err := b.db.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(logsName).Cursor()
	if first, _ := cursor.First(); first == nil {
		return 0, nil
	} else {
		return bytesToUint64(first), nil
	}
}

func (b *Store) GetLog(idx uint64, log *Log) error {
	tx, err := b.db.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	bucket := tx.Bucket(logsName)
	v := bucket.Get(uint64ToByte(idx))
	if v == nil {
		return KeyNotFoundError
	}
	return decodeMsgPack(v, log)
}

func (b *Store) StoreLog(log *Log) error {
	return b.StoreLogs([]*Log{log})
}

func (b *Store) StoreLogs(logs []*Log) error {
	tx, err := b.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, log := range logs {
		k := uint64ToByte(log.Index)
		v, err := encodeMsgPack(log)
		if err != nil {
			return err
		}
		bucket := tx.Bucket(logsName)
		if err := bucket.Put(k, v.Bytes()); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (b *Store) DeleteRange(min, max uint64) error {
	minKey := uint64ToByte(min)
	tx, err := b.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	curs := tx.Bucket(logsName).Cursor()

	for k, _ := curs.Seek(minKey); k != nil; k, _ = curs.Next() {
		if bytesToUint64(k) > max {
			break
		}
		if err := curs.Delete(); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (b *Store) Sync() error {
	return b.db.Sync()
}

func uint64ToByte(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func decodeMsgPack(byte []byte, out interface{}) error {
	r := bytes.NewBuffer(byte)
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(r, &hd)
	return dec.Decode(out)
}

func encodeMsgPack(in interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	hd := codec.MsgpackHandle{}
	enc := codec.NewEncoder(buf, &hd)
	err := enc.Encode(in)
	return buf, err
}
