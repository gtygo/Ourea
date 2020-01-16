package client

import (
	"context"

	pb "github.com/gtygo/Ourea/rpc/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StartClient(addr string) pb.CrudClient {
	logrus.Info("[RPC] client started")
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logrus.Fatalln("[RPC] client got error:", err)
	}
	defer conn.Close()

	return pb.NewCrudClient(conn)
}

func Set(k string, v string, c pb.CrudClient, ctx context.Context) error {
	req := &pb.SetRequest{
		Key:   k,
		Value: v,
	}
	r, err := c.Set(ctx, req)
	if err != nil {
		logrus.Fatal("[RPC] client set error:", err)
	}
	logrus.Infof("[RPC] client received %s", r.Message)
	return nil
}

func Get(k string, c pb.CrudClient, ctx context.Context) (string, error) {
	reqGet := &pb.GetRequest{
		Key: k,
	}
	r, err := c.Get(ctx, reqGet)
	if err != nil {
		logrus.Fatal("[RPC] client get error", err)
		return r.Message, err
	}
	logrus.Infof("[RPC] client received %s , %s", r.Message, r.Value)
	return r.Value, nil
}

func Del(k string, c pb.CrudClient, ctx context.Context) error {
	reqDel := &pb.DelRequest{
		Key: k,
	}
	r, err := c.Delete(ctx, reqDel)
	if err != nil {
		logrus.Fatalf("[RPC] client delete error: %s", err)
		return err
	}
	logrus.Infof("[RPC] client received %s", r.Message)
	return nil
}
