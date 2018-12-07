package main

import (
	"net"
)

type Guochaner interface {
	Guochan(conn net.Conn) net.Conn
}

type listener struct {
	net.Listener
	Guochaner
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	return l.Guochan(c), err
}

func Listen(g Guochaner, network, address string) (net.Listener, error) {
	l, err := net.Listen(network, address)
	return &listener{l, g}, err
}

func Dial(g Guochaner, network, address string) (net.Conn, error) {
	c, err := net.Dial(network, address)
	return g.Guochan(c), err
}

type Method uint

const (
	AES_128_CTR Method = 1 + iota
	AES_192_CTR
	AES_256_CTR
	maxMethod
)

var keySizes = []uint8{
	AES_128_CTR: 16,
	AES_192_CTR: 24,
	AES_256_CTR: 32,
}

/*
var methods = make([]func() (Guochaner, error), maxMethod)

methods[AES_128_CTR] = NewCTRConn
methods[AES_192_CTR] = NewCTRConn
methods[AES_256_CTR] = NewCTRConn

func ( m Method) New() net.Conn {
	if m > 0 && m < maxMethod {
		f := methods[m]
		if f != nil {
			return f()
		}
	}

	pannic("Guochan: requested method function #" + strconv.Itoa(int(m)) + "is unavailable")
}

func ( m Method ) Size() int {
	if h > 0 && h < maxMethod {
			return int(keySizes[m])
	}
	panic("Guochan: Size of unknown method function")
}
*/

func (m Method) Guochan(c net.Conn) net.Conn {
	return NewConn(m, c)
}
