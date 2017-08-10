package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

func handler(conn net.Conn) {
	var buf = bufio.NewReader(conn)
	line, _ := buf.ReadString('\n')
	//	fmt.Println(line)
	line = strings.TrimRightFunc(line, check)
	fmt.Println(line)
	uk := strings.Split(line, " ")
	fmt.Println(uk)
	if uk[0] != "123321" {
		conn.Close()
		return
	}

	uu, _ := ioutil.ReadFile("/root/cfg.sql")
	n2, err := conn.Write(uu)
	if err == nil {
		fmt.Println(n2)
		conn.Close()
	} else {
		fmt.Println("写出错")
	}
}

func check(c rune) bool {
	if c == '\r' || c == '\n' {
		return true
	}
	return false
}

func main() {
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
