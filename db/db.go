package db

import (
	"github.com/gtygo/Ourea/boltkv"
	"github.com/gtygo/Ourea/kv"
)

type DB struct {
	Item kv.Item
}

func NewDB() DB {
	dbs, _ := boltkv.Open()
	s := boltkv.BoltItem{
		Db: dbs,
	}
	return DB{
		Item: &s,
	}
}
