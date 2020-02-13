package main

import (
	"github.com/gtygo/Ourea/server"
	"github.com/gtygo/Ourea/db"

)

func main() {
	server := &server.Server{
		RedisPort: ":3306",
		DbName:  "my.db",
		Db:db.NewDB("my.db"),
	}
	server.StartServer()
}
