package main

import "strconv"

var minParamCount = map[string]int{
	"quit":     0,
	"ping":     0,
	"set":      2,
	"get":      1,
	"del":      1,
	"incr":     2,
	"decr":     2,
	"transfer": 3,
	"pay":      2,
	"invoice":  3,
}

func getMsat(arg []byte) (int64, error) {
	return strconv.ParseInt(string(arg), 10, 64)
}
