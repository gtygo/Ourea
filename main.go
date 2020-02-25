package main

import (
	"flag"

	"github.com/gtygo/Ourea/db"
	"github.com/gtygo/Ourea/raft/node"
	rpc "github.com/gtygo/Ourea/raft/rpc/server"
	"github.com/gtygo/Ourea/server"
)

var RpcAddr string
var RedisAddr string
var raftId int
var dbName string


func init(){
	flag.StringVar(&RpcAddr,"p","127.0.0.1:5001","rpc address")
	flag.StringVar(&RedisAddr,"r",":3306","redis port")
	flag.StringVar(&dbName,"m","db1","data base name")
	flag.IntVar (&raftId,"id",1,"raft id")

}


func main() {

	flag.Parse()
	n := node.NewNode(RpcAddr, uint64(raftId), db.NewDB(dbName))
	svr := &server.Server{
		RedisPort: RedisAddr,
		Rpc:       rpc.NewRpcServer(n),
	}

	go svr.StartKVServer()
	go svr.Rpc.StartRpcServer()
	n.Loop()

}
