package main

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

func ListenAndServe(network, address string) error {

	l, err := net.Listen(network, address)

	if err != nil {
		return err
	}

	for {

		conn, err := l.Accept()

		if err != nil {
			log.Println(err)
			return errors.New("Accept error")
		}

		go HandleConn(conn)

	}

}

func HandleConn(conn net.Conn) {

	err := Negotiate(conn)

	if err != nil {
		log.Println(err)
		return
	}

	err = HandleRequest(conn)

	if err != nil {
		log.Println(err)
		return
	}

	dst, err := GetAddr(conn)

	if err != nil {
		log.Println(err)
		return
	}

	dstconn, err := net.Dial("tcp", dst.String())

	if err != nil {
		err = SendReply(conn, 0x03)

		if err != nil {
			log.Println(err)
		}

		return
	}
	err = SendReply(conn, StatusSucceeded)

	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	f := func(dst, src net.Conn) {
		defer wg.Done()
		io.Copy(dst, src)
		dst.Close()
	}
	go f(conn, dstconn)
	go f(dstconn, conn)
	log.Println("proxy:", conn.RemoteAddr(), "<->", dst.String(), "(", dstconn.RemoteAddr(), ")")

	wg.Wait()

}
