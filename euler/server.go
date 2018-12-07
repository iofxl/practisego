package main

import (
	"fmt"
	"io/ioutil"
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
			log.Fatal(err)
		}

		go func() {

			defer conn.Close()

			sb, err := ioutil.ReadAll(conn)
			fmt.Println(string(sb))

			if err != nil {
				log.Fatal(err)
			}

		}()
	}

}
