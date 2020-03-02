package server

import "errors"

var ErrSetArgs =errors.New("(error) ERR wrong number of arguments for 'set' command")

var ErrGetArgs =errors.New("(error) ERR wrong number of arguments for 'get' command")

var ErrDelArgs =errors.New("(error) ERR wrong number of arguments for 'del' command")

var ErrKeysArgs =errors.New("(error) ERR wrong number of arguments for 'keys' command")

var ErrHsetArgs =errors.New("(error) ERR wrong number of arguments for 'hset' command")

var ErrHgetArgs =errors.New("(error) ERR wrong number of arguments for 'hget' command")

var ErrIncrArgs =errors.New("(error) ERR wrong number of arguments for 'incr' command")

var ErrDecrArgs =errors.New("(error) ERR wrong number of arguments for 'decr' command")

var ErrSaveArgs =errors.New("(error) ERR wrong number of arguments for 'save' command")


