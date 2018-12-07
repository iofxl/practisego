package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {

	r := flag.String("r", "127.0.0.1", "remote ip")
	p := flag.String("p", "12345", "remote port")
	n := flag.String("n", "tcp", "network")

	flag.Parse()

	hostport := net.JoinHostPort(*r, *p)

	laddr, err := net.ResolveTCPAddr(*n, hostport)

	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen(laddr.Network(), laddr.String())

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
