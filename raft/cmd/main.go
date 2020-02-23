package main

import (
	"fmt"
	"github.com/gtygo/raft/rpc/client"
)

func main(){
	client:=client.Client{
		ReqVoteCh:           make(chan uint64, 10),
		ReqVoteDoneCh:       make(chan struct{}, 1),
		AppendEntriesCh:     make(chan uint64, 10),
		AppendEntriesRespCh: make(chan uint64, 10),
	}

	resp,err:=client.DoClientCommand("127.0.0.1:5001",nil)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println("收到leader 的回复：",resp)
}