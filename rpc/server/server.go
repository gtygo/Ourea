package server

import (
	"context"
	"log"
	"net"

	"github.com/gtygo/Ourea/kv"
	pb "github.com/gtygo/Ourea/rpc/pb"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func StartServer(item kv.Item) {
	println("start server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	s := grpc.NewServer()
	pb.RegisterCrudServer(s, &server{
		item: item,
	})
	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}

type server struct {
	item kv.Item
}

func (s *server) Get(context.Context, *pb.GetRequest) (*pb.GetReply, error) {
	panic("implement me")
}

func (s *server) Delete(context.Context, *pb.DelRequest) (*pb.DelReply, error) {
	panic("implement me")
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	log.Println("received: ", in.Key, in.Value)
	if err := s.item.Set([]byte(in.Key), []byte(in.Value)); err != nil {
		return &pb.SetReply{
			Message: err.Error(),
		}, err
	}
	return &pb.SetReply{
		Message: "set success! " + in.Key + "---" + in.Value,
	}, nil
}
