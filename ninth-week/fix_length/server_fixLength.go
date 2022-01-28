package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

const (
	DefinedLength = 1024
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		var msg = make([]byte, DefinedLength)
		n, err := reader.Read(msg)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("Decode message failed, err:", err)
			return
		}
		fmt.Println("Received message from client, the message is:", string(msg[:n]))
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
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
