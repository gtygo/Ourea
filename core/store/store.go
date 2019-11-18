package store

import (
	"github.com/gtygo/Ourea/core/snapshot"
	"github.com/hashicorp/raft"
	"log"
	"os"
)

type Store struct {
	RaftDir  string
	RaftBind string

	raft *raft.Raft

	// Generate snapshot, restore snapshot.
	fsm *snapshot.Fsm

	logger *log.Logger
}

func NewStore(raftDir, raftBind string) (*Store, error) {
	fsm, err := snapshot.NewFsm(raftDir)
	if err != nil {
		return nil, err
	}
	return &Store{
		RaftDir:  raftDir,
		RaftBind: raftBind,
		fsm:      fsm,
		logger:   log.New(os.Stderr, "[store] ", log.LstdFlags),
	}, nil
}
