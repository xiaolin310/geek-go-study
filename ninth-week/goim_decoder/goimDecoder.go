package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Operation uint32

type Package struct {
	Length        uint32
	HeaderLength  uint16
	Version       uint16
	Operation     Operation
	SequenceID    uint32
	Body          []byte
}

const (
	Auth      Operation = 7
	HeartBeat Operation = 2
	Message   Operation = 4
)

func decodeLength(p *Package, r io.Reader) error {
	lengths := make([]byte, 6)
	n, err := r.Read(lengths)
	if err != nil {
		return err
	}
	if n < len(lengths) {
		return errors.New("cannot read enough length data")
	}
	p.Length = decodeUnit32(lengths)
	p.HeaderLength = decodeUnit16(lengths[4:])
	return nil

}

func decodeHeader(p *Package, r io.Reader) error {
	header := make([]byte, p.HeaderLength-6)
	n, err := r.Read(header)
	if err != nil {
		return err
	}
	if n < len(header) {
		return errors.New("cannot read enough header data")
	}
	p.Version = decodeUnit16(header)
	p.Operation = Operation(decodeUnit32(header[2:]))
	p.SequenceID = decodeUnit32(header[6:])
	return nil
}


func decodeUnit16(b []byte) uint16 {

	return binary.BigEndian.Uint16(b)

}

func decodeUnit32(b []byte) uint32 {

	return binary.BigEndian.Uint32(b)
}

// Decode : 只负责解码，不处理EOF error
func Decode(r io.Reader) error {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	p := Package{}
	err := decodeLength(&p, br)
	if err != nil {
		return err
	}
	err = decodeHeader(&p, br)
	if err != nil {
		return err
	}
	p.Body = make([]byte, p.Length-uint32(p.HeaderLength))
	n, err := br.Read(p.Body)
	if err != nil {
		return err
	}
	if n < len(p.Body) {
		return errors.New("cannot read enough body data")
	}
	fmt.Printf("The goim package decoded as:\npackageLen: %d, headerLen: %d, version %d," +
		" operation: %d, sequenceId: %d, body: %v\n", p.Length, p.HeaderLength, p.Version,
		 p.Operation, p.SequenceID, string(p.Body))
	return nil

}
