package main

import (
	"fmt"
	"log"
	"net"
	"os"
	//	"time"
)

func main() {

	// host:port
	address := os.Args[1]

	laddr, err := net.ResolveTCPAddr("tcp4", address)

	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenTCP(laddr.Network(), laddr)

	if err != nil {
		log.Fatal(err)
	}

	for {

		conn, err := l.Accept()

		if err != nil {
			panic(err)
		}

		go echo(conn)

	}

}

func echo(conn net.Conn) {

	defer conn.Close()

	var b [15]byte

	for {

		n, err := conn.Read(b[0:])

		if err != nil {
			return
		}
		fmt.Println(n, err, b[:n], string(b[:n]))
		_, err = conn.Write(b[:n])

		if err != nil {
			panic(err)
		}
	}

}
