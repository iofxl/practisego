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

	return nil

}

func HandleConn(conn net.Conn) {

	err := Negotiate(conn)

	if err != nil {
		log.Println("Negotiation error", err)
		return
	}

	err = HandleRequest(conn)

	if err != nil {
		log.Println("HandleRequset error", err)
		return
	}
	/*

		dst := make([]byte, 2048)

		n, err := conn.Read(dst)

		if err != nil {
			log.Println(err)
			return
		}
	*/

	srvconn, err := net.Dial("tcp", "127.0.0.1:12345")

	if err != nil {
		err = SendReply(conn, 0x03)

		if err != nil {
			log.Println(err)
		}

		return
	}

	/*
		_, err = srvconn.Write(dst[:n])

		if err != nil {

			log.Println(err)
			return
		}
	*/

	var wg sync.WaitGroup
	wg.Add(2)
	f := func(dst, src net.Conn) {
		defer wg.Done()
		io.Copy(dst, src)
		dst.Close()
	}
	go f(conn, srvconn)
	go f(srvconn, conn)

	wg.Wait()

}
