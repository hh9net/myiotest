package main

import (
	"fmt"
	"io"
	"net"
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
	var wg sync.WaitGroup
	cc := make(chan string, 10000)
	go func() {
		for {
			fmt.Println(<-cc)
		}
	}()
	for i := 1; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Add(-1)
			}()
			for {
				conn, err := net.Dial("tcp", "192.168.106.35:10020")
				if err != nil {
					fmt.Println("连接服务端失败:", err.Error())
					return
				}
				//var lengthBytes []byte = make([]byte, 4)
				//if t, err := io.ReadFull(conn, lengthBytes); err != nil {
				conn.Write([]byte("get /data/upfiles/images/2017-01/21/72_scrollpic_new_14849733141.png\n"))

				//buf := bufpool.Get().([]byte)
				//fmt.Println(len(buf))
				//defer bufpool.Put(buf)
				//		var all []byte
				for {
					buf := bufpool.Get().([]byte)
					nr, er := conn.Read(buf)
					if nr > 0 {
						//		nw, ew := dst.Write(buf[0:nr])
						//		fmt.Println(string(buf[0:nr]), len(buf), nr)
						//				all = append(all, buf[0:nr]...)
						fmt.Println(nr)
					}
					if er == io.EOF {
						break
					}
					if er != nil {
						err = er
						break
					}
					bufpool.Put(buf)
					//			fmt.Println(len(all))
					//time.Sleep(time.Second / 100)
				}
				cc <- "搞定一个"
				//	fmt.Println(string(all), len(all))
				//fd, _ := os.OpenFile("ump.png", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
				//	fd_content := strings.Join([]string{"======", fd_time, "=====", str_content, "\n"}, "")
				//	fd.Write(all)
				//	fd.Close()
				conn.Close()
			}
		}()
	}
	wg.Wait()
}
