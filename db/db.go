package db

import (
	"github.com/gtygo/Ourea/boltkv"
	"github.com/gtygo/Ourea/kv"
)

//DB maintains an item interface for implementing a database instance
type DB struct {
	Item kv.Item
}

//NewDB return a DB instance
func NewDB(name string) kv.Item {
	dbs, _ := boltkv.Open(name)
	s := boltkv.BoltItem{
		Db: dbs,
	}
	return &s
}
