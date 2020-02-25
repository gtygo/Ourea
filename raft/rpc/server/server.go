package rpcserver

import (
	"context"
	"net"

	"github.com/gtygo/Ourea/raft/node"
	"github.com/gtygo/Ourea/raft/rpc/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedRaftServiceServer

	node *node.Node
}

func (s *Server) RequestVote(ctx context.Context, req *pb.RequestVoteReq) (*pb.RequestVoteResp, error) {
	logrus.Info("收到了请求投票rpc..... ",req,s.node)

	s.node.Client.ReqVoteCh <- 0

	if req.Term < s.node.CurrentTerm {
		logrus.Infof("request vote path 1")
		return &pb.RequestVoteResp{
			Term:        s.node.CurrentTerm,
			VoteGranted: false,
		}, nil
	}

	if s.node.VotedFor == 0 || s.node.VotedFor == req.CandidateId {
		logrus.Infof("request vote path 2")
		s.node.CurrentTerm=req.Term
		//todo: 检查候选人日志是否和自己同样新
		return &pb.RequestVoteResp{
			Term:        s.node.CurrentTerm,
			VoteGranted: true,
		}, nil
	}
	logrus.Infof("request vote path 3")
	return &pb.RequestVoteResp{
		Term:        req.Term,
		VoteGranted: false,
	}, nil
}

func (s *Server) AppendEntries(ctx context.Context, req *pb.AppendEntriesReq) (*pb.AppendEntriesResp, error) {
	logrus.Info("收到附加日志rpc..... ",req,s.node)
	s.node.Client.AppendEntriesCh<-0

	if len(req.Entries)!=0{
		//not heart beat
		logrus.Infof("rpc server received append entries")
		s.node.HandleAppendEntriesInfo(req.Entries,s.node.CurrentTerm)
	}

	return &pb.AppendEntriesResp{
		Term:    10,
		Success: true,
	}, nil
}


func (s *Server)HandleCommand(data []string){



	if len(data)<=2{
		data=append(data,data[1])
	}
	s.HandleClientCommand(nil,&pb.ClientCommandReq{
		Ins:                  &pb.Instruction{
			Type:                 data[0],
			Key:                  data[1],
			Value:                data[2],
		},
	})
}

func (s *Server)HandleClientCommand(ctx context.Context,req *pb.ClientCommandReq) (*pb.ClientCommandResp, error){




	logrus.Info("收到客户端请求rpc..... ")

	s.node.ClientReqCh<-*req.Ins
	logrus.Info("等待leader处理完成......")
	<- s.node.ClientReqDoneCh
	logrus.Infof("逻辑层处理客户端请求已完成")
	s.node.HandleAppendEntriesInfo([]*pb.Instruction{req.Ins},s.node.CurrentTerm)


	return &pb.ClientCommandResp{Success:true}, nil
}


func NewRpcServer(n *node.Node) *Server {
	return &Server{
		node:   n,
	}
}

func (s *Server)StartRpcServer() {
	logrus.Infof("raft rpc server start listening at: %s ...", s.node.Addr)
	lis, err := net.Listen("tcp", s.node.Addr)
	if err != nil {
		logrus.Fatal(err)
	}
	grpcs := grpc.NewServer()
	pb.RegisterRaftServiceServer(grpcs, s)
	if err := grpcs.Serve(lis); err != nil {
		logrus.Fatal(err)
	}
}

func (s *Server)IsLeader()bool{
	return s.node.NodeState==node.LEADER
}

func (s *Server)Get(k []byte)([]byte,error){
	return s.node.Db.Get(k)
}

func (s *Server)Set(k []byte,v []byte)error{
	return s.node.Db.Set(k,v)
}

func (s *Server)Del(k []byte)error{
	return s.node.Db.Delete(k)
}

