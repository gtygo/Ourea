package snapshot

import (
	"github.com/gtygo/Ourea/core/store"
	"log"
	"os"
)

type Fsm struct {
	db     store.DB
	logger *log.Logger
}

func NewFsm(path string) (*Fsm, error) {
	db, err := NewBadgerDB(path, path)
	if err != nil {
		return nil, err
	}
	return &Fsm{
		db:     db,
		logger: log.New(os.Stderr, "[fsm] ", log.LstdFlags),
	}, nil
}
