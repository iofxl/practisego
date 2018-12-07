package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	address := flag.String("s", "127.0.0.1:12345", "address")
	network := flag.String("p", "tcp", "udp,tcp,both")
	flag.Parse()

	if *network == "tcp" {
		daytimesTCP(*address)
	} else if *network == "udp" {
		daytimesUDP(*address)

	} else if *network == "both" {
		go daytimesTCP(*address)
		daytimesUDP(*address)
	} else {
		flag.Usage()
	}

}

func daytimesTCP(address string) {
	ln, err := net.Listen("tcp", address)

	defer ln.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {

		conn, err := ln.Accept()

		// defer won't work here
		// defer conn.Close()

		if err != nil {
			panic(err)
		}

		daytime := time.Now().String()

		//fmt.Fprintf(conn, "%s\n", []byte(daytime))

		_, err = conn.Write([]byte(daytime))

		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintln(conn)

		conn.Close()

	}

}

func daytimesUDP(address string) {

	pkconn, err := net.ListenPacket("udp", address)

	defer pkconn.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {
		p := make([]byte, 512)
		// 和TCP不一样，这个必须读一下，才知道有人连过来，再写转
		_, addr, err := pkconn.ReadFrom(p)
		if err != nil {
			return
		}
		daytime := time.Now().String()
		_, err = pkconn.WriteTo([]byte(daytime), addr)

		if err != nil {
			return
		}

	}

}
