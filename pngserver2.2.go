package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"strings"
	"sync"
)

var bufpool *sync.Pool

func init() {
	bufpool = &sync.Pool{}
	bufpool.New = func() interface{} {
		return make([]byte, 32*1024)
	}
}
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	buf := bufpool.Get().([]byte)
	defer bufpool.Put(buf)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

var cao = strings.Repeat("xxxxxxxxxxxxxxxx", 1)
var caob = []byte(cao)
var cc int

func handler(conn net.Conn) {
	var buf = bufio.NewReader(conn)
	line, _ := buf.ReadString('\n')
	//	fmt.Println(line)
	line = strings.TrimRightFunc(line, check)
	fmt.Println(line)
	//	uk := strings.Split(line, " ")
	//	fmt.Println(uk)
	//	if uk[0] != "123321" {
	//		conn.Close()
	//		return
	//	}
	//	fd, _ := os.OpenFile("/root/cfg.sql", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	//	Copy(conn, fd)
	//	fd.Close()
	//cc++
	//fmt.Println(cc)
	conn.Write(caob)
	conn.Close()

}

func check(c rune) bool {
	if c == '\r' || c == '\n' {
		return true
	}
	return false
}

func main() {
	runtime.GOMAXPROCS(4)
	ln, err := net.Listen("tcp", "0.0.0.0:10020")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("conn err:", err)
		}
		//	id := time.Now().UnixNano()
		go handler(conn)
	}
}
