package main

import (
	"log"
	"net"
	"testing"
)

func TestDelimiter(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:30001")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for i := 0; i < 10; i++ {
		msg := []byte("Hello, send msg split by delimiter!\n")
		conn.Write(msg)
	}
}