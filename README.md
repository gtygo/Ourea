# Ourea
[![Build Status](https://travis-ci.org/gtygo/Ourea.svg?branch=master)](https://travis-ci.org/gtygo/Ourea)
[![Go Report Card](https://goreportcard.com/badge/github.com/gtygo/Ourea)](https://goreportcard.com/report/github.com/gtygo/Ourea)

High performance K-V cache server

## features
* Supports redis protocol
* Map based on segmented lock,More efficient concurrent reading and writing

## start 
```
go build -o ourea 


./ourea
```

