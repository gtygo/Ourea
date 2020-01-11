package main

import (
	"github.com/gtygo/Ourea/db"
	"github.com/gtygo/Ourea/rpc/server"
)

func main() {
	opt := db.Opt{
		Name: "db1.db",
		Mode: 0666,
	}
	newDB := db.NewDB(opt)
	defer newDB.Item.Close()
	newDB.Item.Set([]byte("123112"), []byte("qwer"))
	server.StartServer()

}
