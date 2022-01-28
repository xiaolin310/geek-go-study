package main

import (
	"fmt"
	"log"
	"net"
	"testing"
)

func TestLengthField(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:30003")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for i := 0; i < 10; i++ {
		msg := "send message encoded by length field based frame!"
		data, err := Encode(msg)
		if err != nil {
			fmt.Println("Encode msg failed, err:", err)
			return
		}
		conn.Write(data)
	}



}
