package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		err := Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("Decode message failed, err:", err)
			return
		}
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:30010")
	if err != nil {
		fmt.Println("Listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept failed, err:", err)
			continue
		}
		go handleConn(conn)
	}
}
