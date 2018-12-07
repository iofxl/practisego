package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func ListenAndServeC(cfg *Config) {

	network := "tcp4"
	address := cfg.Address

	l, err := net.Listen(network, address)

	if err != nil {
		log.Fatalln("Listen error")
	}

	for {

		conn, err := l.Accept()

		if err != nil {
			log.Println("Accept error", err)
		}

		go handleConn(cfg, conn)
	}

}

func handleConn(cfg *Config, conn net.Conn) {

	g := cfg.M
	server := cfg.Server

	err := Negotiate(conn)

	if err != nil {
		log.Println("Negotiation error:", err)
	}

	err = HandleRequest(conn)
	if err != nil {
		log.Println("HandleRequest error:", err)
		return
	}

	b := make([]byte, 2048)

	n, err := conn.Read(b)

	if err != nil {
		log.Printf("GetAddr error:", err)
		return
	}

	dst, err := ParseAddr(b[:n])

	if err != nil {
		log.Println("ParseAddr error:", err)
		return
	}

	srvconn, err := Dial(g, "tcp", server)

	if err != nil {

		SendReply(conn, 0x01)
		log.Println("DialServer error:", err)
		return
	}

	if _, err := srvconn.Write(b[:n]); err == nil {
		log.Printf("Connect: %s\n", dst.String())
	} else {
		log.Printf("Write dst %s error: %v\n", dst.String(), err)
	}

	reply, _ := ReadReply(srvconn)
	SendReply(conn, reply)

	var wg sync.WaitGroup
	f := func(dst, src net.Conn) {
		wg.Add(1)
		defer wg.Done()
		io.Copy(dst, src)
		dst.Close()
	}
	go f(conn, srvconn)
	go f(srvconn, conn)
	wg.Wait()

}
