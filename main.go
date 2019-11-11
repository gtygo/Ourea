package main
import (
	"flag"
	"fmt"
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
}
