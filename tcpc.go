package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

// func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
func main() {

	if len(os.Args) != 2 {
		fmt.Println("wrong")
		os.Exit(1)
	}

	address := os.Args[1]

	raddr, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		log.Fatal(err)
	}

	tcpconn, err := net.DialTCP("tcp", nil, raddr)

	if err != nil {
		log.Fatal(err)
	}

	_, err = tcpconn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))

	if err != nil {
		log.Fatal(err)
	}

	resp, err := ioutil.ReadAll(tcpconn)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp))
}
