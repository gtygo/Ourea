package node

import (
	"math/rand"
	"time"

	"github.com/gtygo/raft/rpc/client"
	"github.com/gtygo/raft/rpc/pb"
	"github.com/sirupsen/logrus"
)

const (
	FOLLOWER = iota
	CANDIDATE
	LEADER
)

type State int

type LogInfo struct {
	Ins  pb.Instruction
	Term uint64
}

type Node struct {
	Id              uint64
	Addr            string
	Time            time.Duration
	NodeState       State
	voteCount       uint64
	ClusterNodeAddr []string
	Client          client.Client
	ClientReqCh     chan pb.Instruction
	ClientReqDoneCh chan struct{}
	LogRespCount int
	StateMachine map[string]string


	//Persistent state on all servers:
	CurrentTerm uint64
	VotedFor    uint64
	Log         []LogInfo

	//Volatile state on all servers:
	CommitIndex uint64
	LastApplied uint64

	//Volatile state on leaders:
	NextIndex  []uint64
	MatchIndex []uint64
}

func NewNode(addr string,id uint64) *Node {

	n := &Node{
		Id:              id,
		Addr:            addr,
		NodeState:       FOLLOWER,
		VotedFor:        0,
		LogRespCount:    0,
		ClientReqCh:make(chan pb.Instruction,10),
		ClientReqDoneCh:make(chan struct{},1),
		StateMachine:make(map[string]string,10),
		ClusterNodeAddr: []string{"127.0.0.1:5001", "127.0.0.1:5002", "127.0.0.1:5003"},
		Client: client.Client{
			ReqVoteCh:           make(chan uint64, 10),
			ReqVoteDoneCh:       make(chan struct{}, 1),
			AppendEntriesCh:     make(chan uint64, 10),
			AppendEntriesRespCh: make(chan uint64, 10),
		},
	}

	return n
}

func (n *Node) Loop() {

	logrus.Infof("startup raft node: %v ,node state is %v (follower)", n.Id, n.NodeState)

	for {

		switch n.NodeState {

		case FOLLOWER:
			logrus.Infof("now state is follower")

			select {
			case x := <-n.Client.ReqVoteCh:
				logrus.Infof("follower received request vote msg %v", x)
				n.checkCurrentTerm(x)
			case x := <-n.Client.AppendEntriesCh:
				logrus.Infof("follower received append entries msg %v", x)
				n.checkCurrentTerm(x)
			case <-time.After(30 * time.Second):
				logrus.Infof("follower timeout ,start new round,%v", 1)
				//todo: start leader election
				n.startLeaderElection()
			}

		case CANDIDATE:
			logrus.Infof("now state is candidate")
			select {
			case x := <-n.Client.ReqVoteCh:
				logrus.Infof("candidate received request vote msg %v", x)
				n.checkCurrentTerm(x)
				n.voteCount++
				if n.voteCount>uint64(len(n.ClusterNodeAddr)+1)/2{
					n.NodeState=LEADER
				}
			case x := <-n.Client.AppendEntriesCh:
				logrus.Infof("candidate received append entries msg %v", x)
				n.checkCurrentTerm(x)

			case <-time.After(25 * time.Second):
				logrus.Infof("candidate timeout ,start new round,%v", 1)
				//todo: start leader election
				n.startLeaderElection()

			}

		case LEADER:
			logrus.Infof("now state is leader")

			//不断发送心跳给follower节点
			// todo： 接收客户端请求
			select 	{
			//收到了客户端发来的请求
			case x:= <- n.ClientReqCh:
				logrus.Infof("逻辑层收到了客户端请求rpc 信号")
				//发送append entries rpc
				entries:=make([]*pb.Instruction,0)
				entries=append(entries,&x)
				n.boardCastHeartBeat(&pb.AppendEntriesReq{
					Term:                 n.CurrentTerm,
					LeaderId:             n.Id,
					PrevLogIndex:         0,
					PrevLogTerm:          0,
					Entries:              entries,
					LeaderCommit:         0,
				},true)
			//收到appendentries 的回复
			case x:=<-n.Client.AppendEntriesCh:
				logrus.Infof("leader received append entries log resp: %v",x)
				n.LogRespCount++
				//收到follower节点都改变了自身的状态机
				if n.LogRespCount==2{
					n.ClientReqDoneCh<- struct{}{}
				}

			default:
				//发送heart beat
				logrus.Infof("leader send normal heart beat rpc")
				time.Sleep(time.Second*5)
				n.boardCastHeartBeat(&pb.AppendEntriesReq{},false)
			}
		}
	}
}

func (n *Node) checkCurrentTerm(T uint64) {
	if T > n.CurrentTerm {
		n.CurrentTerm = T
		n.NodeState = FOLLOWER
	}
}

func (n *Node) startLeaderElection() {
	logrus.Infof("start leader election......")
	n.LogRespCount=0

	// 1.切换状态为候选人
	n.NodeState = CANDIDATE


	// 2.自增任期号
	n.CurrentTerm++
	// 3.重置超时计数器


	n.voteCount=1
	// 4.发送请求投票rpc
	n.boardCastRequestVote()
}

func (n *Node) switchToCandidate() {

}

func (n *Node) receiveVote() {
	logrus.Infof("start listening request vote rpc message....")
	count := 0
	for {
		select {
		case x := <-n.Client.ReqVoteCh:
			count++
			handleRequestVoteMessage(x)

		case <-n.Client.ReqVoteDoneCh:
			count = 0
			handleRequestVoteDoneMessage()
		}
	}

}

//不断发送心跳给其他服务器
func (n *Node) boardCastHeartBeat( req *pb.AppendEntriesReq,isNeedResp bool) {
	for i := 0; i < len(n.ClusterNodeAddr); i++ {
		if n.ClusterNodeAddr[i] == n.Addr {
			continue
		}
		go n.Client.DoAppendEntries(n.ClusterNodeAddr[i], req,isNeedResp)
	}
}

func (n *Node) boardCastRequestVote() {
	for i := 0; i < len(n.ClusterNodeAddr); i++ {
		if n.ClusterNodeAddr[i] == n.Addr {
			continue
		}
		go n.Client.DoRequestVote(n.ClusterNodeAddr[i], &pb.RequestVoteReq{
			Term:         n.CurrentTerm,
			CandidateId:  n.Id,
			LastLogIndex: 0,
			LastLogTerm:  0,
		})
	}
}

func handleRequestVoteMessage(term uint64) {

}

func handleRequestVoteDoneMessage() {

}

func handleAppendEntriesMessage(term uint64) {

}

func (n *Node)HandleAppendEntriesInfo(entries []*pb.Instruction,term uint64){

	for _,x:=range entries {
		//log append
		n.Log = append(n.Log, LogInfo{
			Ins:  *x,
			Term: term,
		})
		//change state machine

		if x.Type=="SET"{
			n.StateMachine[x.Key]=x.Value
		} else if x.Type=="DEL"{
			delete(n.StateMachine,x.Key)
		}
	}

	logrus.Infof("==== now state machine is :")
	logrus.Info("logs: ",n.Log)
	logrus.Info(n.StateMachine)

}

func Random(l int, r int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(r-l) + r
}
