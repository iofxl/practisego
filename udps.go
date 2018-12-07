package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	server := os.Args[1]

	laddr, err := net.ResolveUDPAddr("udp", server)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP(laddr.Network(), laddr)

	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {
		b := make([]byte, 1024)

		_, raddr, err := conn.ReadFromUDP(b)
		_, err = conn.WriteToUDP(b, raddr)

		if err != nil {
			log.Fatal(err)
		}

		b = nil

	}
	fmt.Println()
}
