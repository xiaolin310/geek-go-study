package main

import (
	"encoding/binary"
	"log"
	"net"
	"testing"
)

func encode(msg string) []byte {
	headerLen := 16
	// Auth, HeartBeat, Message
	operation := HeartBeat
	version, sequence := 1, 10000
	packageLen := headerLen + len(msg)
	pkg := make([]byte, packageLen)
	binary.BigEndian.PutUint32(pkg[:4], uint32(packageLen))
	binary.BigEndian.PutUint16(pkg[4:6], uint16(headerLen))
	binary.BigEndian.PutUint16(pkg[6:8], uint16(version))
	binary.BigEndian.PutUint32(pkg[8:12], uint32(operation))
	binary.BigEndian.PutUint32(pkg[12:16], uint32(sequence))
	copy(pkg[16:], []byte(msg))
	return pkg

}

func TestDecode(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:30010")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for i := 0; i < 3; i++ {
		msg := "Test goim Decoder..."
		data := encode(msg)
		conn.Write(data)
	}
}
