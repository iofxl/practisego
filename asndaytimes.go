package main

import (
	"encoding/asn1"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	address := "127.0.0.1:12345"

	l, err := net.Listen("tcp4", address)

	if err != nil {
		log.Fatal(err)
	}

	for {

		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		daytime := time.Now()

		mdata, err := asn1.Marshal(daytime)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(mdata, daytime)
		_, err = conn.Write(mdata)

		if err != nil {
			log.Fatal(err)
		}

		conn.Close()

	}
}
