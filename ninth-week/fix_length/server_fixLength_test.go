package main

import (
	"log"
	"net"
	"testing"
)

func TestFixLength(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for i := 0; i < 10; i++ {
		msg := []byte("Hello, send fix length of package!")
		fixLength := DefinedLength - len(msg)
		msg = append(msg, make([]byte, fixLength)...)

		conn.Write(msg)
	}
}
