package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {

	address := flag.String("h", "127.0.0.1", "address")
	data := flag.String("d", "", "data")
	flag.Parse()

	conn, err := net.Dial("ip4:icmp", *address)

	if err != nil {
		log.Fatal(err)
	}

	msg := []byte{8, 0, 0, 0, 0, 13, 0, 37}

	msg = append(msg, []byte(*data)...)

	resp := make([]byte, len(msg)+20)

	chksum := checkSum(msg)

	// 取高8位和取低8位
	msg[2] = byte(chksum >> 8)
	msg[3] = byte(chksum & 0xff)

	for {

		n, err := conn.Write(msg)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Msg Write:", msg[:n])

		n, err = conn.Read(resp)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Msg Read:", resp[20:n])
	}

}

func checkSum(msg []byte) uint16 {

	var sum int32
	l := len(msg)

	// 先补全成偶数个字节
	if l%2 != 0 {
		msg = append(msg, 0)
	}

	// 当 i == l -2 , i + 1 < l 还是可以成立, 而 i+=2 后 i == l 不再执行循环
	for i := 0; i < l; i += 2 {
		sum += int32(msg[i])<<8 + int32(msg[i+1])
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)

	return uint16(^sum)

}
