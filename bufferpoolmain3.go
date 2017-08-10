package main

import (
	"fmt"
	"io"
	"net"
	"os"
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
func main() {
	conn, err := net.Dial("tcp", "192.168.106.35:10020")
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		return
	}
	//var lengthBytes []byte = make([]byte, 4)
	//if t, err := io.ReadFull(conn, lengthBytes); err != nil {
	conn.Write([]byte("123321 get /data/upfiles/images/2017-01/21/72_scrollpic_new_14849733141.png\n"))
	fd, _ := os.OpenFile("ump.sql", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	//	fd_content := strings.Join([]string{"======", fd_time, "=====", str_content, "\n"}, "")
	//fd.Write([]byte(all.String()))
	//Copy(fd, conn)
	Copy(conn, fd)
	fd.Close()
	conn.Close()
}
