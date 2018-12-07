package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {

	network := flag.String("n", "tcp", "protocol:udp,tcp,both")
	address := flag.String("a", "127.0.0.1:12345", "address")
	flag.Parse()

	if *network == "tcp" {
		echosTCP(*address)
	} else if *network == "udp" {
		echosUDP(*address)
	} else if *network == "both" {
		go echosTCP(*address)
		echosUDP(*address)
	} else {
		flag.Usage()
	}

}

func echosTCP(address string) {

	ln, err := net.Listen("tcp", address)

	defer ln.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		tr := io.TeeReader(conn, os.Stdout)

		// 弗知道这种直接Copy，跟借一个[]byte，先读再写有什么区别？
		_, err = io.Copy(conn, tr)

		if err != nil {
			log.Fatal(err)
		}

	}
}

func echosUDP(address string) {

	pkconn, err := net.ListenPacket("udp", address)

	defer pkconn.Close()

	if err != nil {
		log.Fatal(err)
	}

	p := make([]byte, 512)

	for {
		n, addr, err := pkconn.ReadFrom(p)

		if err != nil {
			panic(err)
		}

		fmt.Println(n, err, p[:n], string(p[:n]))

		_, err = pkconn.WriteTo(p[:n], addr)

		if err != nil {
			panic(err)
		}

	}

}
