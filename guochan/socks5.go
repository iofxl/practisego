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
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"practisego/guochan/guochan"
)

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
	// ErrVersion is ErrVersion
	ErrVersion = errors.New("guochan: version error")
	// ErrCmdNotSupported is ErrCmdNotSupported
	ErrCmdNotSupported = errors.New("guochan: command not supported")
	// ErrATYPNotSupported is ErrATYPNotSupported
	ErrATYPNotSupported = errors.New("guochan: address type not supported")
)

// An Addr represents a SOCKS-specific address.
// Either Name or IP is used exclusively.
type Addr struct {
	Name string // fully-qualified domain name
	IP   net.IP
	Port int
}

// Network is Network
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

// Socks5Server is ...
type Socks5Server struct {
	Addr     string
	Hostname string
	cfg      Config
	g        guochan.G国产器
	logger   *log.Logger
}

// Serve is Serve
func (s *Socks5Server) Serve(l net.Listener) error {

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		ss := s.newSession(c)
		go ss.Serve()

	}

}

// ListenAndServe is ListenAndServe
func (s *Socks5Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":1080"
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(l)

}

// ListenAndServe is ListenAndServe
func ListenAndServe(addr string, cfg Config, logger *log.Logger) error {

	g := guochan.New国产器(guochan.Method(cfg.Method), []byte(cfg.Secret))

	s := &Socks5Server{Addr: addr, cfg: cfg, g: g, logger: logger}

	return s.ListenAndServe()
}

// Session is Session
type Session struct {
	s *Socks5Server
	net.Conn
}

func (s *Socks5Server) newSession(conn net.Conn) *Session {
	return &Session{
		s:    s,
		Conn: conn,
	}
}

// Serve is ...
func (ss *Session) Serve() {
	defer ss.Close()

	err := ss.Negotiate()
	if err != nil {
		ss.s.logger.Println(err)
		return
	}

	// read version
	ver, err := ss.ReadByte()
	if err != nil {
		ss.s.logger.Println(err)
		return
	}

	if ver != Version5 {
		ss.s.logger.Println(ErrVersion)
		return
	}

	// read CMD
	cmd, err := ss.ReadByte()

	if err != nil {
		ss.s.logger.Println(err)
		return
	}

	switch cmd {
	case CmdConnect:
		// read rsv
		_, err := ss.ReadByte()
		if err != nil {
			ss.s.logger.Println(err)
			return
		}

		addr, rawaddr, err := ss.ReadAddrBytes()
		if err != nil {
			ss.s.logger.Println(err)
			return
		}

		midconn, err := guochan.DialTimeout(ss.s.g, "tcp", ss.s.cfg.Server, 3*time.Second)

		if err != nil {
			ss.s.logger.Println(err)
			return
		}
		_, err = midconn.Write(rawaddr)
		if err != nil {
			ss.s.logger.Println(err)
			return
		}

		ss.s.logger.Printf("proxy: %s", addr.String())
		Proxy(ss, midconn)

	default:
		ss.s.logger.Println(ErrCmdNotSupported)
	}

}

// ReadByte is ReadByte
func (ss *Session) ReadByte() (byte, error) {
	b := make([]byte, 1, 1)
	_, err := ss.Read(b)
	return b[0], err
}

// Negotiate is Negotiate
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

	_, err = ss.Read(methods)

	if err != nil {
		return err
	}

	_, err = ss.Write([]byte{Version5, AuthMethodNotRequired})

	return err

}

// ReadAddrBytes is used for read dst addr
func (ss *Session) ReadAddrBytes() (Addr, []byte, error) {

	var addr Addr
	var bufaddr []byte

	atyp, err := ss.ReadByte()
	bufaddr = append(bufaddr, atyp)
	if err != nil {
		return addr, bufaddr, err
	}

	switch atyp {

	case AddrTypeIPv4:
		b := make([]byte, 6)
		n, err := ss.Read(b)

		bufaddr = append(bufaddr, b[:n]...)
		if err != nil {
			return addr, bufaddr, err
		}
		addr.IP = net.IP(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])

	case AddrTypeIPv6:
		b := make([]byte, 18)
		n, err := ss.Read(b)
		bufaddr = append(bufaddr, b[:n]...)
		if err != nil {
			return addr, bufaddr, err
		}
		addr.IP = net.IP(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])

	case AddrTypeFQDN:
		// read url length
		urllen, err := ss.ReadByte()
		bufaddr = append(bufaddr, urllen)
		if err != nil {
			return addr, bufaddr, err
		}
		b := make([]byte, urllen+2)
		n, err := ss.Read(b)
		bufaddr = append(bufaddr, b[:n]...)
		if err != nil {
			return addr, bufaddr, err
		}
		addr.Name = string(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])

	default:
		return addr, bufaddr, ErrATYPNotSupported
	}
	return addr, bufaddr, nil
}

// Proxy is ...
func Proxy(l, r net.Conn) {

	var wg sync.WaitGroup

	f := func(dst, src net.Conn) {
		defer wg.Done()
		_, err := io.Copy(dst, src)
		if err != nil {
			dst.Close()
		}
	}

	wg.Add(2)
	go f(l, r)
	go f(r, l)

	wg.Wait()

}
