package core

import (
	"github.com/hashicorp/raft"
	"log"
	"os"
)

type Store struct {
	RaftDir  string
	RaftBind string

	raft *raft.Raft

	// Generate snapshot, restore snapshot.
	fsm *fsm

	logger *log.Logger
}

func NewStore(raftDir, raftBind string) (*Store, error) {
	fsm, err := NewFsm(raftDir)
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

func (s *Store) Open(isStrap bool, localID string) error {

	return nil
}

func (s *Store) Get(key string) (string, error) {
	return s.fsm.Get(key)
}

func (s *Store) Set(key, value string) error {
	return s.fsm.Set(key, value)
}

func (s *Store) Delete(key string) error {
	return nil
}

func (s *Store) Join(id, addr string) error {
	return nil
}

func (s *Store) Leave(id string) error {
	return nil
}

func (s *Store) SnapShot() error {
	return nil
}
