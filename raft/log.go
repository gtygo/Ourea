package raft

type LogType uint8

type Log struct {
	Index uint64

	Term uint64

	Type LogType

	Data []byte

	Extensions []byte
}

type LogStore interface {
	FirstIndex() (uint64, error)
	LastIndex() (uint64, error)
	GetLog(index uint64, log *Log) error
	StoreLog(log *Log) error
	StoreLogs(logs []*Log) error
	DeleteRange(min, max uint64) error
}
