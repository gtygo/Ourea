package server

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gtygo/Ourea/boltkv"
	"github.com/gtygo/Ourea/cache"
	"github.com/gtygo/Ourea/redis"
	"github.com/sirupsen/logrus"
)

type Server struct {
	RedisPort string
	c *cache.Cache
	OldT map[string]time.Time
}

func NewServer(port string)*Server{
	return &Server{
		RedisPort: port,
		c:cache.NewCache(),
	}
}

func (s *Server) StartKVServer() {
	logrus.Infof("kv server is starting listening at %v \n",s.RedisPort)
	c, err := net.Listen("tcp", s.RedisPort)
	if err!=nil{
		logrus.Fatal("start listener failed !",err)
	}
	//cover dump.rdb
	if err:=s.c.Cover();err!=nil&&err!=boltkv.ErrBucketNotFound{
		logrus.Fatal("cover dump.rdb failed:",err)
	}
	go s.CleanOldKey()
	logrus.Info("rdb file cover success,data length:",len(s.c.GetAllKey()))
	for {
		a, _ := c.Accept()
		go s.handleConn(a)
	}
}

func readAll(c net.Conn) []byte {
	var result []byte
	var buf [2048]byte
	var n = 2048
	var err error
	for n == 2048 {
		n, err = c.Read(buf[:])
		if err != nil {
			return nil
		}
		result = append(result, buf[0:n]...)
	}
	return result
}

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()

	for {
		buf := readAll(c)
		if len(buf) == 0 {
			return
		}
		fmt.Println(string(buf))
		if len(buf) != 0 {
			v, _ := redis.GetReply(buf)
			strData := fmt.Sprintf("%s", v)
			data := handleRedisStr(strData)
			logrus.Info("kv server收到客户端请求 ..... ")

					ans, err := s.DispatchCommand(data)
					if err != nil {
						ans = []string{err.Error()}
					}
					redisResp := redis.GetRequest(append([]string{}, ans...))
					c.Write(redisResp)

			}

	}
}

func handleRedisStr(info string) []string {
	ans := strings.Replace(info, "[", "", 1)
	ans = strings.Replace(ans, "]", "", 1)
	ansArr:=strings.Split(ans, " ")
	ansArr[0]=strings.ToUpper(ansArr[0])
	return ansArr
}

func (s *Server)CleanOldKey(){
	serverStartTime:=time.Now()

	for {
		time.Sleep(time.Second*10)

		for k,_:=range s.OldT{
			if time.Now().Sub(serverStartTime).Seconds() > float64(s.OldT[k].Second()){
				s.c.Del(k)
				delete(s.OldT,k)
			}
		}
	}



}