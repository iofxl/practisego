// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package socks provides a SOCKS version 5 client implementation.
//
// SOCKS protocol version 5 is defined in RFC 1928.
// Username/Password authentication for SOCKS version 5 is defined in
// RFC 1929.
package 国产

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
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

func Negotiate(conn net.Conn) error {

	b := make([]byte, 215)

	_, err := conn.Read(b)

	if err != nil {
		log.Println(err)
		return err
	}

	if b[0] != Version5 {
		return errors.New("Only Support Version5")
	}

	_, err = conn.Write([]byte{Version5, AuthMethodNotRequired})

	return err

}

func HandleRequest(conn net.Conn) error {

	b := make([]byte, 3)

	_, err := io.ReadFull(conn, b)

	if err != nil {

		log.Println(err)
		return err
	}

	if !bytes.Equal([]byte{Version5, CmdConnect, RSV}, b) {
		return errors.New("Only Support CmdConnect")
	}

	return nil

}

func ParseAddr(b []byte) (Addr, error) {

	n := len(b)

	var addr Addr

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

func SendReply(conn net.Conn, reply byte) error {

	_, err := conn.Write([]byte{Version5, reply, RSV, AddrTypeIPv4, 0, 0, 0, 0, 0, 0})

	return err
}

func ReadReply(conn net.Conn) (uint8, error) {

	b := make([]byte, 10)
	if _, err := io.ReadFull(conn, b); err != nil {
		return 0x01, err
	}

	return uint8(b[1]), nil
}
