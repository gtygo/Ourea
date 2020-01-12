package main

import (
	"github.com/gtygo/Ourea/db"
	"github.com/gtygo/Ourea/rpc/server"
)

func main() {

	newDB := db.NewDB("my.db")
	defer newDB.Item.Close()
	newDB.Item.Set([]byte("123112"), []byte("qwer"))
	server.StartServer(newDB.Item)

}
