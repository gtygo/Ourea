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
	f.logger.Printf("get key: %s from fsm store",key)
	v, err := f.db.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (f *fsm)Set(key ,value string)error{
	f.logger.Printf("set key: %s vale: %s from fsm store",key,value)
	err:=f.db.Set([]byte(key),[]byte(value))
	if err!=nil{
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

func (f *fsm) SnapShot() (hraft.FSMSnapshot, error) {
	return nil, nil
}

func (f *fsm) Restore(readClose io.ReadCloser) error {
	return nil
}
