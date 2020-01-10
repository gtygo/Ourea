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

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := pb.NewCrudClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Set(ctx, &pb.SetRequest{
		Key:   "key1",
		Value: "value1",
	})
	if err != nil {
		log.Fatal("set error:", err)
	}
	log.Println("get response :", r.GetMessage())

}

func main(){

	StartClient()

}