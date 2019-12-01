package core

import (
	boltraft "github.com/gtygo/Ourea/raft"
	"github.com/hashicorp/raft"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
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
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(localID)
	conf.SnapshotThreshold = 1024
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		s.logger.Println("kmbga")
		return err
	}
	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		s.logger.Println("563hgd")
		return err
	}

	ss, err := raft.NewFileSnapshotStore(s.RaftDir, 2, os.Stderr)
	if err != nil {
		s.logger.Println("9u8y7")
		return err
	}

	boltDB, err := boltraft.NewStore(filepath.Join(s.RaftDir, "raft.db"))
	if err != nil {
		s.logger.Println("43234")
		return err
	}

	r, err := raft.NewRaft(conf, s.fsm, boltDB, boltDB, ss, transport)
	if err != nil {
		s.logger.Println("12 ")
		return err
	}
	s.raft = r
	if isStrap {
		config := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      conf.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		s.raft.BootstrapCluster(config)
	}

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
