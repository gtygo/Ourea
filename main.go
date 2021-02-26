package main

import "github.com/gtygo/Ourea/server"

func main(){
	
        s:=server.NewServer(":3306")
	s.StartKVServer()
}
