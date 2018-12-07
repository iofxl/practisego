package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func ListenAndServeS(network, address string) error {

	l, err := net.Listen(network, address)

	if err != nil {
		return err
	}

	for {

		conn, err := l.Accept()

		if err != nil {
			log.Println(err)
		}

		go handleConnS(conn)

	}

}

func handleConnS(conn net.Conn) {

	dst, err := GetAddr(conn)

	if err != nil {

		return
	}

	dstconn, err := net.Dial("tcp", dst.String())

	if err != nil {
		log.Println(err)
		err = SendReply(conn, 0x03)
		if err != nil {
			log.Println(err)
		}
		return
	}

	SendReply(conn, 0x00)

	var wg sync.WaitGroup
	wg.Add(2)
	f := func(dst, src net.Conn) {
		defer wg.Done()
		io.Copy(dst, src)
		dst.Close()
	}
	go f(conn, dstconn)
	go f(dstconn, conn)
	log.Printf("proxy: %s <-> %s <-> %s(%s)\n", conn.RemoteAddr(), conn.LocalAddr(), dst.String(), dstconn.RemoteAddr())

	wg.Wait()

}
