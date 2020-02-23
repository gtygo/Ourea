package main

import (
	"flag"
	"fmt"
	"github.com/gtygo/raft/rpc/server"
	"strconv"

	"github.com/gtygo/raft/node"
)

var DefaultAddr string
var Id string

func init() {
	flag.StringVar(&DefaultAddr, "addr", "127.0.0.1:5001", "rpc ip address+port")
	flag.StringVar(&Id,"id","1","raft node id")
}

func main() {

	flag.Parse()

	fmt.Println(DefaultAddr,Id)
	id,_:=strconv.Atoi(Id)
	n := node.NewNode(DefaultAddr,uint64(id))
	go server.StartRpcServer(n)
	n.Loop()
}
