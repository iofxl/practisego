package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	server := os.Args[1]

	laddr, err := net.ResolveTCPAddr("tcp", server)

	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenTCP(laddr.Network(), laddr)

	for {

		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		handleclient := func(conn net.Conn) {

			defer conn.Close()
			count := 0

			br := bufio.NewReader(conn)

			for {

				line, err := br.ReadBytes('\n')

				if err != nil {
					log.Fatal(err)
				}

				conn.Write(line)
				count++

				fmt.Println(count)
			}

		}

		go handleclient(conn)
	}

}
