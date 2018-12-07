package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {

	client := flag.Bool("c", false, "client")
	flag.Parse()

	address := "127.0.0.1:12345"

	if *client {
		base64client(address)
	} else {
		base64server(address)
	}

}

func base64server(address string) {

	l, err := net.Listen("tcp4", address)

	defer l.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {

		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		handleconn(conn)

	}

}

func handleconn(conn net.Conn) {

	buf := new(bytes.Buffer)

	tr := io.TeeReader(conn, buf)
	r := base64.NewDecoder(base64.StdEncoding, tr)
	wc := base64.NewEncoder(base64.StdEncoding, conn)

	// 总是忘记，这个无法 var b []byte这样的，这样的len(b) = cap(b) = 0
	b := make([]byte, 512)

	n, err := r.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("base64:", buf.String())
	buf.Reset()
	fmt.Println("len:", n, "data:", b[:n], string(b[:n]))

	_, err = wc.Write(b[:n])

	if err != nil {
		log.Fatal(err)
	}

	err = wc.Close()

	if err != nil {
		log.Fatal(err)
	}

}

func base64client(address string) {

	conn, err := net.Dial("tcp4", address)

	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	r := base64.NewDecoder(base64.StdEncoding, conn)

	var content []byte
	b := make([]byte, 512)

	fmt.Scanln(&content)

	fmt.Println("input:", string(content))
	wc := base64.NewEncoder(base64.StdEncoding, conn)
	n, err := wc.Write(content)

	if err != nil {
		log.Fatal(err)
	}

	wc.Close()

	n, err = r.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("len:", n, "data:", b[:n], string(b[:n]))

}
