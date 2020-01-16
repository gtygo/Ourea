package main

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func CmdRoot(addr string) {
	// c:=client.StartClient(addr)
	logrus.Info("[CMD] start waiting for command line...")
	for {
		print("> ")
		rd := bufio.NewReader(os.Stdin)
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				println("bye")
				return
			}
			logrus.Warnf("[CMD] client got error: %s when read command line")
			println("got error ,pls retry")
		}
		println(line[:len(line)-1])
	}
}

func dispatch() {
}

func parser(line string) {
	println("parse input command line ....")
}
