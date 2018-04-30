package main

import (
	"testing"
	"net"
	"time"
	"encoding/binary"
	"strings"
)

func Test_toBytes(t *testing.T) {
	correct := []byte{0, 0, 0, 1, 0, 2}
	result := toBytes([]uint16{0, 1, 2})

	if len(result) != len(correct) {
		t.Errorf("Fail! Different lengths")
	}

	for i := range result {
		if correct[i] != result[i] {
			t.Errorf("Fail! Incorrect result")
		}
	}
}

func Test_sendError(t *testing.T) {
	conn := &mockConn{}
	var remoteAddress net.UDPAddr
	message := "Testing"

	sendError(conn, &remoteAddress, message)

	/* Should receive byte array
	   2 bytes opcode = ERROR
	   2 bytes errorcode = 0
	   n bytes message
	   1 byte = 0 */
	opCode := binary.BigEndian.Uint16(conn.message[0:2])
	errorCode := binary.BigEndian.Uint16(conn.message[2:4])

	// Note: Stripping any extra 0s from the end of the received string
	received := strings.Trim(string(conn.message[4:]), string(0))

	if opCode != ERROR {
		t.Errorf("Fail! Wrong op code returned")
	}

	if errorCode != 0 {
		t.Errorf("Fail! Error code should be 0 (not defined)")
	}

	if message != received {
		t.Errorf("Fail! Didn't receive the correct message")
	}
}

func Test_sendAck(t *testing.T) {
	conn := &mockConn{}
	var remoteAddress net.UDPAddr
	var testBlockNumber uint16 = 2

	sendAck(conn, &remoteAddress, testBlockNumber)

	/* Should receive byte array
	   2 bytes opcode = ACK
	   2 bytes block number */
	opCode := binary.BigEndian.Uint16(conn.message[0:2])
	blockNumber := binary.BigEndian.Uint16(conn.message[2:4])

	if opCode != ACK {
		t.Errorf("Fail! Wrong op code returned")
	}

	if blockNumber != testBlockNumber {
		t.Errorf("Fail! Wrong block number returned")
	}
}

func Test_printFile(t *testing.T) {
	var m map[uint16][]byte
	var result bool

	/* Complete file case */
	m = make(map[uint16][]byte)
	m[1] = []byte("First blob")
	m[2] = []byte("Second Blob")

	result = printFile(m, 2)
	if !result {
		t.Errorf("Fail! Should have been ok!")
	}

	/* Incomplete file case */
	m = make(map[uint16][]byte)
	m[1] = []byte("First blob")
	m[3] = []byte("Third Blob")

	result = printFile(m, 3)
	if result {
		t.Errorf("Fail! Should have been a fail!")
	}
}

/*
Mocked out the PacketConn interface.
https://golang.org/src/net/net.go?s=9837:12035#L294

type PacketConn interface {
	ReadFrom(b []byte) (n int, addr Addr, err error)
	WriteTo(b []byte, addr Addr) (n int, err error)
	Close() error
	LocalAddr() Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}*/

type mockConn struct {
	message []byte
}

func (t *mockConn) WriteTo(message []byte, address net.Addr) (n int, err error) {
	t.message = message
	return 0, nil
}

func (t *mockConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	return 0, nil, nil
}

func (t *mockConn) Close() error {
	return nil
}

func (t *mockConn) LocalAddr() net.Addr {
	return nil
}

func (t *mockConn) SetDeadline(time time.Time) error {
	return nil
}

func (t *mockConn) SetReadDeadline(time time.Time) error {
	return nil
}

func (t *mockConn) SetWriteDeadline(time time.Time) error {
	return nil
}
