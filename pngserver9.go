package main

import (
	"fmt"
	"sync"
)

var bufpool *sync.Pool

func init() {
	bufpool = &sync.Pool{}
	bufpool.New = func() interface{} {
		return make([]byte, 32*1024)
	}
}

func main() {
	buf := bufpool.Get().([]byte)
	buf = []byte("caon")
	fmt.Println(string(buf))
	defer bufpool.Put(buf)
}
