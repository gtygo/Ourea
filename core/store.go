package core

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
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

var (
	ErrNotLeader = errors.New("not leader")
)

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
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)
	config.SnapshotThreshold = 1024

	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	ss, err := raft.NewFileSnapshotStore(s.RaftDir, 2, os.Stderr)
	if err != nil {
		return err
	}

	// boltDB implement log store and stable store interface
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft.db"))
	if err != nil {
		return err
	}

	// raft system
	r, err := raft.NewRaft(config, s.fsm, boltDB, boltDB, ss, transport)
	if err != nil {
		return err
	}
	s.raft = r

	if isStrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		s.raft.BootstrapCluster(configuration)
	}
	return nil
}

func (s *Store) Get(key string) (string, error) {
	return s.fsm.Get(key)
}

func (s *Store) Set(key, value string) error {
	if s.raft.State() != raft.Leader {
		return ErrNotLeader
	}
	c := NewSetCommand(key, value)

	msg, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(msg, 10*time.Second)

	return f.Error()

}

func (s *Store) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return ErrNotLeader
	}

	c := NewDeleteCommand(key)

	msg, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(msg, 10*time.Second)

	return f.Error()
}

func (s *Store) Join(id, addr string) error {
	s.logger.Printf("received join request for remote node %s, addr %s", id, addr)

	cf := s.raft.GetConfiguration()

	if err := cf.Error(); err != nil {
		s.logger.Printf("failed to get raft configuration")
		return err
	}

	for _, server := range cf.Configuration().Servers {
		if server.ID == raft.ServerID(id) {
			s.logger.Printf("node %s already joined raft cluster", id)
			return nil
		}
	}

	f := s.raft.AddVoter(raft.ServerID(id), raft.ServerAddress(addr), 0, 0)
	if err := f.Error(); err != nil {
		return err
	}

	s.logger.Printf("node %s at %s joined successfully", id, addr)

	return nil
}

func (s *Store) Leave(id string) error {
	s.logger.Printf("received leave request for remote node %s", id)

	cf := s.raft.GetConfiguration()

	if err := cf.Error(); err != nil {
		s.logger.Printf("failed to get raft configuration")
		return err
	}

	for _, server := range cf.Configuration().Servers {
		if server.ID == raft.ServerID(id) {
			f := s.raft.RemoveServer(server.ID, 0, 0)
			if err := f.Error(); err != nil {
				s.logger.Printf("failed to remove server %s", id)
				return err
			}

			s.logger.Printf("node %s leaved successfully", id)
			return nil
		}
	}

	s.logger.Printf("node %s not exists in raft group", id)

	return nil
}

func (s *Store) SnapShot() error {
	s.logger.Printf("doing snapshot mannually")
	f := s.raft.Snapshot()
	return f.Error()
}
