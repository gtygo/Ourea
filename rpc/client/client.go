package client

import (
	"bufio"
	"context"
	"errors"
	pb "github.com/gtygo/Ourea/rpc/pb"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	MsgSuccess = "Success"
	MsgFailed  = "Failed"
)

var (
	ErrCommandArgs = errors.New("command args error")
)

func StartClient(addr string) {
	logrus.Info("[RPC] client started")
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logrus.Fatalln("[RPC] client got error:", err)
	}
	defer conn.Close()
	c := pb.NewCrudClient(conn)
	ctx := context.Background()
	for {
		print("> ")
		rd := bufio.NewReader(os.Stdin)
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				println("bye")
				return
			}
			logrus.Warnf("[RPC] client got error: %s when read command line")
			println("got error ,pls retry")
		}
		args := parser(line)
		value, err := dispatch(args, c, ctx)
		if err != nil {
			logrus.Errorf("[RPC] client dispatch got error: %s", err)
		}
		println(value)
		if value == "bye" {
			return
		}

		//println(line[:len(line)-1])
	}

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

func parser(line string) []string {
	str := strings.Join(strings.Fields(line), " ")
	subCmd := strings.Split(str, " ")
	return subCmd
}

func dispatch(args []string, c pb.CrudClient, ctx context.Context) (string, error) {
	if len(args) == 0 {
		return "", ErrCommandArgs
	}
	handleName := strings.ToUpper(args[0])
	switch handleName {
	case "GET":
		if len(args) != 2 {
			return "", ErrCommandArgs
		}
		return Get(args[1], c, ctx)
	case "SET":
		if len(args) != 3 {
			return "", ErrCommandArgs
		}
		return "", Set(args[1], args[2], c, ctx)
	case "DEL":
		if len(args) != 2 {
			return "", ErrCommandArgs
		}
		return "", Del(args[1], c, ctx)

	case "EXIT":
		return "bye", nil
	}

	return "", ErrCommandArgs
}
