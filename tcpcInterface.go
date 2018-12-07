package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {

	address := flag.String("s", "bing.com:80", "address")
	flag.Parse()

	conn, err := net.Dial("tcp", *address)

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))

	if err != nil {
		log.Fatal(err)
	}

	resp, err := ioutil.ReadAll(conn)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp))

}
