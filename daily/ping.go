package main

import (
	"crypto/rand"
	"fmt"
	"log"
	mrand "math/rand"
	"net"
	"os"
	"time"
)

func main() {

	var address string
	if len(os.Args) == 2 {
		address = os.Args[1]
	} else {
		address = "127.0.0.1"
	}

	conn, err := net.Dial("ip4:icmp", address)

	if err != nil {
		log.Fatal(err)
	}

	/*
		msg[0] = 8  // echo
		msg[1] = 0  // code 0
		msg[2] = 0  // checksum, fix later
		msg[3] = 0  // checksum, fix later
		msg[4] = 0  // identifier[0]
		msg[5] = 13 // identifier[1] (arbitrary)
		msg[6] = 0  // sequence[0]
		msg[7] = 37 // sequence[1] (arbitrary)
		len := 8
	*/

	for {
		// this can't put outside for
		msg := []byte{8, 0, 0, 0, 0, 13, 0, 37}

		n := mrand.Intn(1024)

		data := make([]byte, n)

		rand.Read(data)

		msg = append(msg, data...)

		chksum := checkSum(msg)

		// get high 8 and low 8
		msg[2] = byte(chksum >> 8)
		msg[3] = byte(chksum & 0xff)

		n, err := conn.Write(msg)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Msg Len: %v\n", n)
		fmt.Printf("Msg Write:\t%v\n", msg[:8])

		resp := make([]byte, len(msg)+20)

		n, err = conn.Read(resp)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Msg Read:\t%v\n\n", resp[20:20+8])
		time.Sleep(1 * time.Second)
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
