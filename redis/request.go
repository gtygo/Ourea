package redis

import (
	"fmt"
	"strconv"
	"strings"
)

func GetRequest(args []string) []byte {
	req := []string{
		"*" + strconv.Itoa(len(args)),
	}
	for _, arg := range args {
		req = append(req, "$"+strconv.Itoa(len(arg)))
		req = append(req, arg)
	}
	fmt.Println(req)
	str := strings.Join(req, "\r\n")
	return []byte(str + "\r\n")
}
