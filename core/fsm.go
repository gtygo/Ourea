package core

import (
	"io"
	"log"
	"os"

	hraft "github.com/hashicorp/raft"
)

type fsm struct {
	db     DB
	logger *log.Logger
}

//NewFsm return  a fsm instance
func NewFsm(path string) (*fsm, error) {
	db, err := NewDB(path, path)
	if err != nil {
		return nil, err
	}
	return &fsm{
		db:     db,
		logger: log.New(os.Stderr, "[fsm] ", log.LstdFlags),
	}, nil
}

//Get return value
func (f *fsm) Get(key string) (string, error) {
	v, err := f.db.Get([]byte(key))
	if err != nil {
		f.logger.Fatalf("get key %s error: %s", key, err)
		return "", err
	}
	return string(v), nil
}

func (f *fsm) Set(key, value string) error {
	err := f.db.Set([]byte(key), []byte(value))
	if err != nil {
		f.logger.Fatalf("set key: %s value: %s error: %s ", key, value, err)
		return err
	}
	return nil
}

func (f *fsm) Apply(*hraft.Log) interface{} {
	panic("implement me")
}

func (f *fsm) Snapshot() (hraft.FSMSnapshot, error) {
	panic("implement me")
}

func (f *fsm) Restore(readClose io.ReadCloser) error {
	panic("implement")
}
