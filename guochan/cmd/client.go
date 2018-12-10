package cmd

import (
	"io"
	"log"
	"net"
	"practisego/guochan/国产"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use: "client",
	Run: ListenAndServeC,
}

func ListenAndServeC(cmd *cobra.Command, args []string) {

	log.Printf("Listen port: %v\n", cfg.Port)
	log.Printf("Server: %v\n", cfg.Server)
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(cfg.Port))

	if err != nil {
		log.Fatalln("Listen error")
	}

	for {

		conn, err := l.Accept()

		if err != nil {
			log.Println("Accept error", err)
		}

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {

	err := 国产.Negotiate(conn)

	if err != nil {
		log.Println("Negotiation error:", err)
	}

	err = 国产.HandleRequest(conn)
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

	dst, err := 国产.ParseAddr(b[:n])

	if err != nil {
		log.Println("ParseAddr error:", err)
		return
	}

	g := 国产.New国产器(cfg.M, []byte(cfg.Secret))
	srvconn, err := 国产.Dial(g, "tcp", cfg.Server)

	if err != nil {

		国产.SendReply(conn, 0x01)
		log.Println("DialServer error:", err)
		return
	}

	if _, err := srvconn.Write(b[:n]); err == nil {
		log.Printf("Connect: %s\n", dst.String())
	} else {
		log.Printf("Write dst %s error: %v\n", dst.String(), err)
	}

	reply, _ := 国产.ReadReply(srvconn)
	国产.SendReply(conn, reply)

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
