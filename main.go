package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gtygo/Ourea/server"
)

var (
	node string
	port string
)

func init() {
	flag.StringVar(&node, "id", "1", "node id")
	flag.StringVar(&port, "p", "9090", "port number")

}

func main() {
	flag.Parse()
	fmt.Println(node, port)

	server := server.NewServer()
	fmt.Println(server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	<-quit

}
