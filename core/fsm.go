package core

import (
	"io"
	"log"
	"os"

	"github.com/gtygo/Ourea/raft"
	hraft "github.com/hashicorp/raft"
)

type fsm struct {
	db     DB
	logger *log.Logger
}

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

func (f *fsm) Get(key string) (string, error) {
	return "", nil
}

func (f *fsm) Apply(l *raft.Log) interface{} {
	return nil
}

func (f *fsm) SnapShot() (hraft.FSMSnapshot, error) {
	return nil, nil
}

func (f *fsm) Restore(readClose io.ReadCloser) error {
	return nil
}
