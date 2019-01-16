package main

import (
	"flag"
	"log"
	"net/textproto"
	"os"
	"time"

	"github.com/ziutek/telnet"
)

type conn struct {
	*telnet.Conn
}

func main() {

	var addr string
	flag.StringVar(&addr, "s", "", "server")
	flag.Parse()

	t, err := telnet.DialTimeout("tcp", addr, 2*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	c := &conn{t}

	c.cmd("#", "sh running-config")

	buf := make([]byte, 512)

	for {

		n, err := t.Read(buf)

		os.Stdout.Write(buf[:n])

		if err != nil {
			log.Fatal(err)
		}

	}

}

func (c conn) cmd(expect, format string, args ...interface{}) error {

	err := c.SkipUntil(expect)

	if err != nil {
		return err
	}

	text := textproto.NewConn(c)
	return text.PrintfLine(format, args...)
}
