package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	host := os.Args[1]
	port := os.Args[2]

	hp := net.JoinHostPort(host, port)

	fmt.Println(hp)
}
