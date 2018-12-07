package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

const (
	Version5 = 0x05

	AddrTypeIPv4 = 0x01
	AddrTypeFQDN = 0x03
	AddrTypeIPv6 = 0x04

	AuthMethodNotRequired = 0x00
)

type Addr struct {
	Name string
	IP   net.IP
	Port int
}

func main() {

	var network, address string

	flag.StringVar(&network, "n", "tcp", "network")
	flag.StringVar(&address, "l", ":12345", "address")
	flag.Parse()

	l, err := net.Listen(network, address)

	if err != nil {
		log.Fatal(err)
	}

	for {

		conn, _ := l.Accept()

		defer conn.Close()

		go func() {

			sb := make([]byte, 1000)

			conn.Read(sb)

			if sb[0] != Version5 {
				return
			}

			conn.Write([]byte{0x05, 0x00})

			n, _ := conn.Read(sb)

			var a Addr

			switch sb[3] {

			case AddrTypeIPv4:
				a.IP = net.IP(sb[4 : 4+4])
			case AddrTypeIPv6:
				a.IP = net.IP(sb[4 : 4+16])
			case AddrTypeFQDN:
				a.Name = string(sb[5 : 5+int(sb[4])])
			default:
				return
			}

			a.Port = int(sb[n-2])<<8 | int(sb[n-1])

			address := net.JoinHostPort(a.Name, strconv.Itoa(a.Port))
			fmt.Println(address)

			dconn, err := net.Dial("tcp", address)

			if err != nil {
				fmt.Println(err)
				return
			}

			defer dconn.Close()

			conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功

			go io.Copy(conn, dconn)
			io.Copy(dconn, conn)

		}()

	}

}
