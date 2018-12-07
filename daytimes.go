package main

import (
	"log"
	"net"
	"time"
)

func main() {

	address := ":12345"

	laddr, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		log.Fatal(err)
	}

	tln, err := net.ListenTCP("tcp4", laddr)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := tln.Accept()

		if err != nil {
			continue
		}

		time := time.Now().String()

		_, err = conn.Write([]byte(time))

		if err != nil {
			log.Fatal(err)
		}

		conn.Close()

	}
}
