package main

import (
	"github.com/gtygo/Ourea/db"
	"github.com/gtygo/Ourea/raft/node"
	rpc "github.com/gtygo/Ourea/raft/rpc/server"
	"github.com/gtygo/Ourea/server"
)

func main() {
	server := &server.Server{
		RedisPort: ":3306",
		DbName:    "my.db",
		Db:        db.NewDB("my.db"),
	}
	go server.StartServer()

	n := node.NewNode("127.0.0.1:5001", 1)
	n.Loop()

	go rpc.StartRpcServer(n)

}
