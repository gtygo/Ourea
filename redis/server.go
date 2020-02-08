package redis

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)


func StartListen(){
	logrus.Infof("# warning: using the default config. In order to specify a config file use redis-server /path/redis.conf")
	c, _ := net.Listen("tcp", ":3306")

	for {
		a, _ := c.Accept()
		go handleConn(a)
	}
}

func readAll(c net.Conn)[]byte{
	var result []byte
	var buf [2048]byte
	var n=2048
	var err error
	for n==2048{
		n,err=c.Read(buf[:])
		if err!=nil{
			return nil
		}
		result=append(result,buf[0:n]...)
	}
	return result
}

func handleConn(c net.Conn) {
	defer c.Close()


	for {
		//buf := make([]byte, 1024)

		//length,err:=bufio.NewReader(c).Read(buf[:]),todo: this is a bug

		buf:=readAll(c)
		fmt.Println(string(buf))
		if len(buf)!=0{
			v,_:=GetReply(buf)
			fmt.Printf("value: %s",v)
		}

		c.Write(buf)

	}
}
