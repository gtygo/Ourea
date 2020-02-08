package redis

import (
	"strconv"
	"strings"
)

const (
	SimpleStringsFirstByteReply = '+'
	ErrorsFirstByteReply        = '-'
	IntegersFirstByteReply      = ':'
	BulkStringsFirstByteReply   = '$'
	ArraysFirstByteReply        = '*'
	OkReply                     = "OK"
	PongReply                   = "PONG"
)

/*
redis protocol:
eg:

SET a b

byte[]:

*3
$3
SET
$1
a
$1
b

string:

"*3\r\n$3\r\nSET\r\n$1\r\na\r\n$1\r\nb\r\n"



*/

func GetReply(reply []byte) (interface{}, error) {
	replyType := reply[0]

	switch replyType {
	case SimpleStringsFirstByteReply:
		return handleSimpleReply(reply)
	case ErrorsFirstByteReply:
		return handleErrorReply(reply)
	case IntegersFirstByteReply:
		return handleIntegerReply(reply)
	case BulkStringsFirstByteReply:
		return handleBulkReply(reply)
	case ArraysFirstByteReply:
		return handleArraysReply(reply)
	default:
		return nil, nil
	}
}

func handleSimpleReply(reply []byte) (string, error) {
	if len(reply) == 3 && reply[1] == 'O' && reply[2] == 'K' {
		return OkReply, nil
	}
	return string(reply), nil
}

func handleErrorReply(reply []byte) (string, error) {
	return string(reply), nil
}

func handleIntegerReply(reply []byte) (int, error) {
	pos := flagCount('\r', reply)
	result, err := strconv.Atoi(string(reply[:pos]))
	if err != nil {
		return 0, err
	}
	return result, nil
}

func handleBulkReply(reply []byte) (interface{}, error) {
	pos := flagCount('\r', reply)
	pot := 0
	if reply[:pos][0] == '$' {
		pot = 1
	}

	vlen, err := strconv.Atoi(string(reply[pot:pos]))
	if err != nil {
		return nil, err
	}
	if vlen == -1 {
		return nil, nil
	}

	start := pos + 2
	end := start + vlen
	return string(reply[start:end]), nil
}

func handleArraysReply(reply []byte) (interface{}, error) {
	replyStrs := strings.Split(string(reply), "\r\n")
	replylen := len(replyStrs)
	replyStrs = replyStrs[1 : replylen-1]

	r := []interface{}{}
	for i := 0; i < replylen-1; i++ {
		if i%2 == 1 {
			rv := strings.Join([]string{
				replyStrs[i-1],
				replyStrs[i],
			}, "\r\n") + "\r\n"

			value, err := handleBulkReply([]byte(rv))
			if err != nil {
				return nil, err
			}

			r = append(r, value)
		}
	}

	return r, nil
}

func flagCount(flag byte, reply []byte) int {
	count := 0

	for i, _ := range reply {
		if flag == reply[i] {
			break
		}
		count++
	}
	return count
}
