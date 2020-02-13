package server

import (
	"fmt"
	"github.com/gtygo/Ourea/redis"
	"net"
	"strings"

	"github.com/gtygo/Ourea/kv"
	"github.com/sirupsen/logrus"
)

type Server struct {
	RedisPort string

	DbName string

	Db kv.Item
}

func (s *Server) StartServer() {
	logrus.Infof("# warning: using the default config. In order to specify a config file use redis-server /path/redis.conf")
	c, _ := net.Listen("tcp", s.RedisPort)
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
			ans,err:=s.DispatchCommand(data)
			if err!=nil{
				ans=err.Error()
			}
			if ans==""{
				ans="OK"
			}
			redisResp:=redis.GetRequest(append([]string{},ans))
			c.Write(redisResp)
		}
	}
}

func handleRedisStr(info string) []string {
	ans := strings.Replace(info, "[", "", 1)
	ans = strings.Replace(ans, "]", "", 1)
	return strings.Split(ans, " ")
}
