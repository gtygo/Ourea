package server

import (
	"context"
	"net"

	"github.com/gtygo/raft/node"
	"github.com/gtygo/raft/rpc/pb"
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

func (s *Server)HandleClientCommand(ctx context.Context,req *pb.ClientCommandReq) (*pb.ClientCommandResp, error){
	logrus.Info("收到客户端请求rpc..... ")

	s.node.ClientReqCh<-*req.Ins
	logrus.Info("等待leader处理完成......")
	<- s.node.ClientReqDoneCh
	logrus.Infof("逻辑层处理客户端请求已完成")
	s.node.HandleAppendEntriesInfo([]*pb.Instruction{req.Ins},s.node.CurrentTerm)


	return &pb.ClientCommandResp{Success:true}, nil
}

func StartRpcServer(n *node.Node) {
	logrus.Infof("raft rpc server start listening at: %s ...", n.Addr)
	lis, err := net.Listen("tcp", n.Addr)
	if err != nil {
		logrus.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterRaftServiceServer(s, &Server{node:n})
	if err := s.Serve(lis); err != nil {
		logrus.Fatal(err)
	}

}
