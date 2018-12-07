package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {

	address := flag.String("r", "127.0.0.1:12345", "address")
	flag.Parse()

	raddr, err := net.ResolveUDPAddr("udp", *address)

	uc, err := net.DialUDP("udp", nil, raddr)

	if err != nil {
		log.Fatal(err)
	}

	_, err = uc.Write([]byte("anything"))

	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 512)

	n, addr, err := uc.ReadFromUDP(b)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(addr, string(b[:n]))

}
