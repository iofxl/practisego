package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ziutek/telnet"
)

type conn struct {
	*telnet.Conn
}

type zxr10 struct {
	Host, Port, User, Pass, Enpass string
}

func (z *zxr10) fetch() (string, error) {

	/*
		Username:
		Password:
		> en 18
		# terminal length 0
		# sh running-config
		# exit
	*/

	port := "21"

	if z.Port != "" {
		port = z.Port
	}

	addr := net.JoinHostPort(z.Host, port)

	t, err := telnet.DialTimeout("tcp", addr, 2*time.Second)

	if err != nil {
		return "", err
	}

	c := conn{t}

	c.cmd("Username:", z.User)
	c.cmd("Password:", z.Pass)
	c.cmd(">", "en 18")
	c.cmd("Password:", z.Enpass)
	c.cmd("#", "terminal length 0")
	c.cmd("#", "sh running-config")
	cfg, err := t.ReadUntil("#")
	if err != nil {
		return string(cfg), err
	}
	c.cmd("", "exit")

	return string(cfg), nil

}

func main() {

	var host, port, user, pass, enpass string
	flag.StringVar(&host, "h", "", "host")
	flag.StringVar(&port, "p", "23", "port")
	flag.StringVar(&user, "u", "", "user")
	flag.StringVar(&pass, "P", "", "passwd")
	flag.StringVar(&enpass, "e", "", "enable passwd")
	flag.Parse()

	z := &zxr10{host, port, user, pass, enpass}
	cfg, err := z.fetch()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)

}

func (c conn) cmd(expect, format string, args ...interface{}) error {

	if expect != "" {

		err := c.SkipUntil(expect)

		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(c, format, args...)
	if err != nil {
		return err
	}

	_, err = c.Write([]byte("\r\n"))
	return err

}
