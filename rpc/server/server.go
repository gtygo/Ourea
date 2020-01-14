package server

import (
	"context"
	"net"

	"github.com/gtygo/Ourea/kv"
	pb "github.com/gtygo/Ourea/rpc/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	port       = ":50051"
	MsgSuccess = "success"
	MsgFailed  = "failed"
)

func StartServer(item kv.Item) {
	logrus.Infof("[RPC] server start, listening at TCP %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Warnf("failed to listen: %s", err)
	}
	s := grpc.NewServer()
	pb.RegisterCrudServer(s, &server{
		item: item,
	})
	if err = s.Serve(lis); err != nil {
		logrus.Warnf("failed to serve: %s", err)
	}
}

type server struct {
	item kv.Item
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	logrus.Infof("[RPC] server received get: %s ", in.Key)
	value, err := s.item.Get([]byte(in.Key))
	if err != nil {
		return &pb.GetReply{
			Value:   "",
			Message: MsgFailed,
		}, err
	}
	return &pb.GetReply{
		Value:   string(value),
		Message: MsgSuccess,
	}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	logrus.Infof("[RPC] server received delete: %s ", in.Key)
	if err := s.item.Delete([]byte(in.Key)); err != nil {
		return &pb.DelReply{
			Message: MsgFailed,
		}, err
	}
	return &pb.DelReply{
		Message: MsgSuccess,
	}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	logrus.Infof("[RPC] server received set: %s , %s", in.Key, in.Value)
	if err := s.item.Set([]byte(in.Key), []byte(in.Value)); err != nil {
		return &pb.SetReply{
			Message: MsgFailed,
		}, err
	}
	return &pb.SetReply{
		Message: MsgSuccess,
	}, nil
}
