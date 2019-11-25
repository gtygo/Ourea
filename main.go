package main

import (
	"flag"
	"fmt"
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
)

func init() {
	flag.StringVar(&port, "port", "4080", "server listing port")
	flag.StringVar(&addr, "addr", ":19090", "raft bind address")
	flag.StringVar(&path, "path", "./", "data directory")
	flag.StringVar(&id, "id", "", "raft node id")
}

func main() {
	flag.Parse()
	conf := config.NewConfig(port, path, id, "", addr)

	svr := server.NewServer(conf)
	fmt.Println(svr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	go svr.Start()
	<-quit

}
