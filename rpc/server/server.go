package server

import (
	"context"
	"log"
	"net"

	pb "github.com/gtygo/Ourea/rpc/pb"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func StartServer() {
	println("start server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	s := grpc.NewServer()
	pb.RegisterCrudServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}

type server struct {
	pb.UnimplementedCrudServer
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	log.Println("received: ", in.Key, in.Value)

	return &pb.SetReply{
		Message: "set success! " + in.Key + "---" + in.Value,
	}, nil
}
