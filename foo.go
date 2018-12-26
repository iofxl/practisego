package main

import (
	"log"
	"net"
)

func main() {

	l, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Print(err)
			return
		}

		go handleConn(conn)

	}
}

func handleConn(conn net.Conn) {
	for {
		b := conn.readByte()
		switch b {
		case msg1:
			print(1)
		case msg2:
			print(2)
		case msg3:
			print(3)
		}
	}
}
