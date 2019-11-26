package server

import (
	"fmt"
	"github.com/gtygo/Ourea/config"
	"github.com/gtygo/Ourea/core"
	"github.com/siddontang/goredis"
	"log"
	"net"
	"os"
)

type Server struct {
	listener net.Listener

	logger *log.Logger

	//TODO : db
	store *core.Store
}

//启动ouera后台服务端，进行一些初始化配置
func NewServer(config *config.Config) *Server {
	server := &Server{
		logger: log.New(os.Stderr, "[server] ", log.LstdFlags),
	}
	listener, err := net.Listen("tcp", config.Port)
	if err != nil {
		panic(err.Error())
	}
	server.listener = listener
	server.logger.Printf("server listening in %s", config.Port)

	store, err := core.NewStore(config.Path, config.Addr)
	if err != nil {
		panic(err.Error())
	}
	server.store = store
	server.logger.Printf("server init store node ,DB path: %s ,link address: %s", config.Path, config.Addr)

	//加入新节点
	if config.Join != "" {
		redisClient := goredis.NewClient(config.Join, "")
		server.logger.Printf("request join to %s", config.Join)
		_, err := redisClient.Do("join", config.Path, config.Id)
		if err != nil {
			server.logger.Println(err)
		}
		redisClient.Close()
	}

	return server
}

func (server *Server) Start() {
	for {
		select {
		default:
			conn, err := server.listener.Accept()
			if err != nil {
				continue
			}
			fmt.Println(conn)
		}

	}
}
