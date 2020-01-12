package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"

	pb "github.com/gtygo/Ourea/rpc/pb"
)

const (
	addr = "localhost:50051"
)

func StartClient() {
	log.Printf("[RPC] client start")
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := pb.NewCrudClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.SetRequest{
		Key:   "key1",
		Value: "value1",
	}
	r, err := c.Set(ctx, req)
	if err != nil {
		log.Fatal("set error:", err)
	}
	log.Printf("[RPC] client received %s", r.Message)

}

func main() {

	StartClient()

}
