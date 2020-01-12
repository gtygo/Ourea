package main

import (
	"github.com/gtygo/Ourea/db"
	"github.com/gtygo/Ourea/rpc/server"
)

func main() {

	newDB := db.NewDB("my.db")
	defer newDB.Item.Close()
	server.StartServer(newDB.Item)

}
