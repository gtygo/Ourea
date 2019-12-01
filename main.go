package main

import (
	"flag"
	"github.com/gtygo/Ourea/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/gtygo/Ourea/server"
)

var (
	port string
	path string
	id   string
	addr string
	join string
)

func init() {
	flag.StringVar(&port, "port", "127.0.0.1:5379", "server listing port")
	flag.StringVar(&addr, "addr", ":19090", "raft bind address")
	flag.StringVar(&path, "path", "./data/", "data directory")
	flag.StringVar(&id, "id", "4235", "raft node id")
	flag.StringVar(&join, "join", "", "join to cluster")
}

func main() {
	flag.Parse()
	conf := config.NewConfig(port, path, id, join, addr)

	svr := server.NewServer(conf)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	go svr.Start()
	<-quit

}
