// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package socks provides a SOCKS version 5 client implementation.
//
// SOCKS protocol version 5 is defined in RFC 1928.
// Username/Password authentication for SOCKS version 5 is defined in
// RFC 1929.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

func main() {

	if err := ListenAndServe(":12345"); err != nil {
		log.Fatal(err)
	}
}

// A Command represents a SOCKS command.
type Command uint8

func (cmd Command) String() string {
	switch cmd {
	case CmdConnect:
		return "socks connect"
	case cmdBind:
		return "socks bind"
	default:
		return "socks " + strconv.Itoa(int(cmd))
	}
}

// An AuthMethod represents a SOCKS authentication method.
type AuthMethod uint8

// A Reply represents a SOCKS command reply code.
type Reply uint8

func (code Reply) String() string {
	switch code {
	case StatusSucceeded:
		return "succeeded"
	case 0x01:
		return "general SOCKS server failure"
	case 0x02:
		return "connection not allowed by ruleset"
	case 0x03:
		return "network unreachable"
	case 0x04:
		return "host unreachable"
	case 0x05:
		return "connection refused"
	case 0x06:
		return "TTL expired"
	case 0x07:
		return "command not supported"
	case 0x08:
		return "address type not supported"
	default:
		return "unknown code: " + strconv.Itoa(int(code))
	}
}

// Wire protocol constants.
const (
	Version5 = 0x05

	RSV = 0x00

	AddrTypeIPv4 = 0x01
	AddrTypeFQDN = 0x03
	AddrTypeIPv6 = 0x04

	CmdConnect = 0x01 // establishes an active-open forward proxy connection
	cmdBind    = 0x02 // establishes a passive-open forward proxy connection

	AuthMethodNotRequired         = 0x00 // no authentication required
	AuthMethodUsernamePassword    = 0x02 // use username/password
	AuthMethodNoAcceptableMethods = 0xff // no acceptable authentication methods

	StatusSucceeded = 0x00
)

var (
	ErrVersion          = errors.New("guochan: only support socks5")
	ErrCmdNotSupported  = errors.New("guochan: command not supported")
	ErrATYPNotSupported = errors.New("guochan: address type not supported")
)

// An Addr represents a SOCKS-specific address.
// Either Name or IP is used exclusively.
type Addr struct {
	Name string // fully-qualified domain name
	IP   net.IP
	Port int
}

func (a *Addr) Network() string { return "socks" }

func (a *Addr) String() string {
	if a == nil {
		return "<nil>"
	}
	port := strconv.Itoa(a.Port)
	if a.IP == nil {
		return net.JoinHostPort(a.Name, port)
	}
	return net.JoinHostPort(a.IP.String(), port)
}

type Server struct {
	Addr     string
	Hostname string
}

func (s *Server) Serve(l net.Listener) error {

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		ss := s.newSession(c)

		go func() {
			err := ss.Serve()
			if err != nil {
				log.Println(err)
				return
			}

		}()

	}
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":1080"
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("Listen on %s\n", addr)
	return s.Serve(l)

}

func ListenAndServe(addr string) error {
	s := &Server{Addr: addr}
	return s.ListenAndServe()
}

type Session struct {
	s    *Server
	conn net.Conn
}

func (s *Server) newSession(conn net.Conn) *Session {
	return &Session{
		s:    s,
		conn: conn,
	}
}

func (ss *Session) Serve() error {
	defer ss.conn.Close()

	err := ss.Negotiate()
	if err != nil {
		return err
	}

	// read version
	ver, err := ss.ReadByte()
	if err != nil {
		return err
	}

	if ver != Version5 {
		return ErrVersion
	}

	// read CMD
	cmd, err := ss.ReadByte()

	if err != nil {
		return err
	}

	switch cmd {
	case CmdConnect:
		// read rsv
		_, err := ss.ReadByte()
		if err != nil {
			return err
		}

		// READ DST
		addr, err := ss.ReadAddr()
		if err != nil {
			return err
		}

		log.Println(addr.String())

		// Dial DST
		dstconn, err := net.Dial("tcp", addr.String())
		if err != nil {
			return err
		}

		// SEND REPLY
		err = ss.SendReply(StatusSucceeded)
		if err != nil {
			return err
		}

		Proxy(ss.conn, dstconn)

	default:
		ss.conn.Close()
		return ErrCmdNotSupported
	}

	return nil

}

func (ss *Session) ReadByte() (byte, error) {
	b := make([]byte, 1, 1)
	_, err := io.ReadFull(ss.conn, b)
	return b[0], err
}

func (ss *Session) Negotiate() error {

	// read version
	ver, err := ss.ReadByte()
	if err != nil {
		return err
	}

	if ver != Version5 {
		return ErrVersion
	}

	// read nmethods
	nmethods, err := ss.ReadByte()

	if err != nil {
		return err
	}

	methods := make([]byte, int(nmethods))

	_, err = ss.conn.Read(methods)

	if err != nil {
		return err
	}

	_, err = ss.conn.Write([]byte{Version5, AuthMethodNotRequired})

	return err

}

func (ss *Session) ReadAddr() (Addr, error) {

	var addr Addr

	atyp, err := ss.ReadByte()
	if err != nil {
		return addr, err
	}

	switch atyp {

	case AddrTypeIPv4:
		b := make([]byte, 6)
		n, err := ss.conn.Read(b)
		if err != nil {
			return addr, err
		}
		addr.IP = net.IP(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])
	case AddrTypeIPv6:
		b := make([]byte, 18)
		n, err := ss.conn.Read(b)
		if err != nil {
			return addr, err
		}
		addr.IP = net.IP(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])
	case AddrTypeFQDN:
		// read url length
		urllen, err := ss.ReadByte()
		if err != nil {
			return addr, err
		}
		b := make([]byte, urllen+2)

		n, err := ss.conn.Read(b)
		if err != nil {
			return addr, err
		}
		addr.Name = string(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])
	default:
		return addr, ErrATYPNotSupported
	}
	return addr, nil
}

func ParseAddr(b []byte) (Addr, error) {

	var addr Addr

	n := len(b)

	addr.Port = int(b[n-2])<<8 | int(b[n-1])

	switch b[0] {
	case AddrTypeIPv4, AddrTypeIPv6:
		addr.IP = net.IP(b[1 : n-2])
	case AddrTypeFQDN:
		addr.Name = string(b[2 : n-2])
	default:
		return addr, errors.New("ParseAddr error")
	}

	return addr, nil

}

func (ss *Session) SendReply(reply byte) error {

	_, err := ss.conn.Write([]byte{Version5, reply, RSV, AddrTypeIPv4, 0, 0, 0, 0, 0, 0})

	return err
}

func ReadReply(conn net.Conn) (uint8, error) {

	b := make([]byte, 10)
	if _, err := io.ReadFull(conn, b); err != nil {
		return 0x01, err
	}

	return uint8(b[1]), nil
}

func Proxy(l, r net.Conn) {

	var wg sync.WaitGroup

	f := func(dst, src net.Conn) {
		defer wg.Done()
		io.Copy(dst, src)
		dst.Close()
	}

	wg.Add(2)
	go f(l, r)
	go f(r, l)
	wg.Wait()

}

func echo() error {
	fmt.Println
}
