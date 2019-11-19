package server

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	listener net.Listener

	logger *log.Logger

	//TODO : db
}

func NewServer() *Server {
	server := &Server{
		listener: nil,
		logger:   log.New(os.Stderr, "[server] ", log.LstdFlags),
	}
	listener, err := net.Listen("tcp", "8090")
	if err != nil {
		fmt.Println(err.Error())
	}
	server.listener = listener
	return server
}

func (server *Server) Start(){
	for {
		select{
		default:
			conn,err:=server.listener.Accept()
			if err!=nil{
				continue
			}
		}

	}
}
