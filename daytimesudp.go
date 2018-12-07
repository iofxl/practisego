package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("wrong")
	}

	address := os.Args[1]

	laddr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		log.Fatal(err)
	}

	uc, err := net.ListenUDP("udp", laddr)

	if err != nil {
		log.Fatal(err)
	}

	for {
		handleudp(uc)
	}
}

func handleudp(udpconn *net.UDPConn) {

	// 这里使b []byte的字节切片就不行,这是为什么呢？
	// 因为没有经过make的切片长度是0，所以就一直被读满，但是err == nil
	// var b []byte
	//var b [512]byte
	b := make([]byte, 5)

	// n, addr, err := udpconn.ReadFromUDP(b[0:])
	n, addr, err := udpconn.ReadFromUDP(b)

	if err != nil {
		return
	}

	fmt.Println(err, b[:n], string(b[:n]))

	daytime := time.Now().String()

	fmt.Println(daytime)

	_, err = udpconn.WriteToUDP([]byte(daytime+"\n"), addr)

	if err != nil {
		return
	}

}
