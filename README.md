# Ourea
[![Build Status](https://travis-ci.org/gtygo/Ourea.svg?branch=master)](https://travis-ci.org/gtygo/Ourea)
[![Go Report Card](https://goreportcard.com/badge/github.com/gtygo/Ourea)](https://goreportcard.com/report/github.com/gtygo/Ourea)

High performance K-V cache server

## features
* Supports redis protocol
* Map based on segmented lock,More efficient concurrent reading and writing
* Supports persistence, similar with RDB of redis

## start
### server 
```
go build -o ourea 

./ourea
```
### client
```
redis-cli -p 3306
```

## support command
* SET 
* GET 
* DEL 
* KEYS 
* HSET 
* HGET
* HDEL
* SAVE
* BGSAVE

## design

![avatar](https://s2.ax1x.com/2020/03/09/89BGUH.png)


## plan
* Support master-slave replication
* Support more data structures
* Support SQL syntax parsing
* Better performance




