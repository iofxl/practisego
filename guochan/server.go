package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func ListenAndServeS(cfg *Config) {

	network := "tcp4"
	address := cfg.Address
	g := cfg.M

	l, err := Listen(g, network, address)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println(err)
			break
		}

		go handleConnS(conn)

	}
}

func handleConnS(conn net.Conn) {

	b := make([]byte, 2048)

	n, err := conn.Read(b)

	if err != nil {
		log.Printf("GetAddr error", err)
		return
	}

	dst, err := ParseAddr(b[:n])

	if err != nil {
		log.Println(err)
		return
	}

	dstconn, err := net.Dial("tcp", dst.String())

	if err != nil {
		log.Println(err)
		SendReply(conn, 0x03)
		return
	}

	SendReply(conn, 0x00)

	var wg sync.WaitGroup
	f := func(dst, src net.Conn) {
		wg.Add(1)
		defer wg.Done()
		io.Copy(dst, src)
		dst.Close()
	}
	go f(conn, dstconn)
	go f(dstconn, conn)
	log.Printf("proxy: %s <-> %s <-> %s(%s)\n", conn.RemoteAddr(), conn.LocalAddr(), dst.String(), dstconn.RemoteAddr())
	wg.Wait()
	log.Println("closed")

}
