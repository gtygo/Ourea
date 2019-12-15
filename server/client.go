package server

import (
	"bufio"
	"errors"
	"github.com/gtygo/Ourea/core"
	"github.com/siddontang/goredis"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	ErrParams        = errors.New("ERR params invalid")
	ErrRespType      = errors.New("ERR resp type invalid")
	ErrCmdNotSupport = errors.New("ERR command not supported")
)

type Client struct {
	server *Server

	store *core.Store

	cmd string

	args [][]byte

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
	defer func(c *Client) {
		c.conn.Close()
	}(c)
	for {
		c.cmd = ""
		c.args = nil

		req, err := c.reader.ParseRequest()
		if err != nil && err != io.EOF {
			c.logger.Println(err.Error())
			return
		} else if err != nil {
			return
		}
		err = c.handleRequest(req)
		if err != nil && err != io.EOF {
			c.logger.Println(err.Error())
			return
		}
	}
}

func (c *Client) handleRequest(req [][]byte) error {
	//todo add error
	if len(req) == 0 {
		c.cmd = ""
		c.args = nil
	} else {
		c.cmd = strings.ToLower(string(req[0]))
		c.args = req[1:]
	}
	var (
		err error
		v   string
	)

	c.logger.Printf("processing %s command", c.cmd)
	//c.logger.Printf("c.cmd: ", c.cmd)
	switch c.cmd {
	case "get":
		if v, err = c.Get(); err == nil {
			c.Resp(v)
		}
	case "set":
		if err = c.Set(); err == nil {
			c.Resp("OK")
		}
	case "del":
		if err = c.Del(); err == nil {
			c.Resp("OK")
		}
	case "join":
		if err = c.Join(); err == nil {
			c.Resp("OK")
		}
	case "leave":
		if err = c.Leave(); err == nil {
			c.Resp("OK")
		}
	case "ping":
		if len(c.args) != 0 {
			err = ErrParams
		}
		c.Resp("PONG")
		err = nil
	case "snapshot":
		if err = c.SnapShot(); err == nil {
			c.Resp("OK")
		}
	default:
		err = ErrCmdNotSupport
	}
	if err != nil {
		c.Resp(err)
	}

	return nil
}
func (c *Client) Get() (string, error) {
	if len(c.args) != 1 {
		return "", ErrParams
	}
	key := string(c.args[0])
	v, err := c.store.Get(key)
	c.logger.Println("get information :", v, err)
	if err != nil {
		c.logger.Printf("got error when get : %s", key)
		return "", err
	}
	return v, nil
}
func (c *Client) Set() error {
	if len(c.args) != 2 {
		return ErrParams
	}
	key := string(c.args[0])
	val := string(c.args[1])
	err := c.store.Set(key, val)
	if err != nil {
		return err

	}
	return nil
}
func (c *Client) Del() error {
	if len(c.args) != 1 {
		return ErrParams
	}
	key := string(c.args[0])
	err := c.store.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) SnapShot() error {
	if len(c.args) != 0 {
		return ErrParams
	}
	err := c.store.SnapShot()
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) Join() error {
	if len(c.args) != 2 {
		return ErrParams
	}
	addr := string(c.args[0])
	id := string(c.args[0])
	return c.store.Join(id, addr)
}
func (c *Client) Leave() error {
	if len(c.args) != 1 {
		return ErrParams
	}
	id := string(c.args[0])
	return c.store.Leave(id)
}

func (c *Client) Resp(v interface{}) error {
	var err error = nil
	switch val := v.(type) {
	case []interface{}:
		err = c.writer.WriteArray(val)
	case []byte:
		err = c.writer.WriteBulk(val)
	case nil:
		err = c.writer.WriteBulk(nil)
	case int64:
		err = c.writer.WriteInteger(val)
	case string:
		err = c.writer.WriteString(val)
	case error:
		err = c.writer.WriteError(val)
	default:
		err = errors.New("resp type error")
	}
	if err != nil {
		return err
	}
	return c.writer.Flush()
}
