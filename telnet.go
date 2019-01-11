package main

import (
	"log"
	"os"

	"github.com/ziutek/telnet"
)

func main() {

	conn, err := telnet.Dial("tcp", "bat.org:23")

	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 256)
	for {

		n, err := conn.Read(buf)
		os.Stdout.Write(buf[:n])
		if err != nil {
			log.Fatal(err)
		}

	}

}
