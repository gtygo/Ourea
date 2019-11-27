package server

import (
	"bufio"
	"bytes"
	"github.com/gtygo/Ourea/core"
	"github.com/siddontang/goredis"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	server *Server

	store *core.Store

	cmd string

	args [][]byte

	buf bytes.Buffer

	conn net.Conn

	reader *goredis.RespReader

	writer *goredis.RespWriter

	logger *log.Logger
}

func InitClient(conn net.Conn, server *Server) {
	c := &Client{
		server: server,
		store:  server.store,
		logger: log.New(os.Stderr, "[client] ", log.LstdFlags),
	}

	c.conn = conn
	c.reader = goredis.NewRespReader(bufio.NewReader(conn))
	c.writer = goredis.NewRespWriter(bufio.NewWriter(conn))

	go c.handleConn()
}

func (c *Client) handleConn() {
	defer c.conn.Close()
	for {
		c.cmd = ""
		c.args = nil

		req, err := c.reader.ParseRequest()
		if err != nil || err != io.EOF {
			c.logger.Println(err.Error())
			return
		}
		err = c.handleRequest(req)
		if err != nil || err != io.EOF {
			c.logger.Println(err.Error())
			return
		}
	}
}

func (c *Client) handleRequest(req [][]byte) error {
	var (
		error error
		//val string
	)

	if len(req) == 0 {
		c.cmd = ""
		c.args = nil
	} else {
		c.cmd = strings.ToLower(string(req[0]))
		c.args = req[1:]
	}

	c.logger.Printf("processing %s command", c.cmd)

	switch c.cmd {
	case "get":
	case "set":
	case "del":
	case "join":
	case "leave":
	case "ping":
	case "snapshot":
	default:

	}
	return error

}
