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
func NewDB() DB {
	dbs, _ := boltkv.Open()
	s := boltkv.BoltItem{
		Db: dbs,
	}
	return DB{
		Item: &s,
	}
}
