package db

import (
	"github.com/gtygo/Ourea/boltkv"
	"github.com/gtygo/Ourea/kv"
	"os"
)

//DB maintains an item interface for implementing a database instance
type DB struct {
	Item kv.Item
}

type Opt struct {
	Name string
	Mode os.FileMode
}

//NewDB return a DB instance
func NewDB(opt Opt) DB {
	dbs, _ := boltkv.Open(opt)
	s := boltkv.BoltItem{
		Db: dbs,
	}
	return DB{
		Item: &s,
	}
}
