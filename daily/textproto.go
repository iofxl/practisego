package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os/exec"
	"strings"
)

func main() {

	var isServer bool
	var host string
	var addr string
	flag.BoolVar(&isServer, "s", false, " is server ")
	flag.StringVar(&host, "h", "127.0.0.1:12345", "host")
	flag.StringVar(&addr, "l", ":12345", "listen addr")
	flag.Parse()

	if isServer {
		if err := ListenAndServe(addr); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := clientFunc(host); err != nil {
			log.Fatal(err)
		}

	}

}

type Client struct {
	Text       *textproto.Conn
	conn       net.Conn
	serverName string
	// map of supported extensions
	ext map[string]string
	// supported auth mechanisms
	auth       []string
	localName  string // the name to use in HELO/EHLO
	didHello   bool   // whether we've said HELO/EHLO
	helloError error  // the error from the hello
}

func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return NewClient(conn, host)
}

func NewClient(conn net.Conn, host string) (*Client, error) {
	text := textproto.NewConn(conn)

	_, _, err := text.ReadResponse(220)

	if err != nil {
		text.Close()
		return nil, err
	}

	c := &Client{
		Text:       text,
		conn:       conn,
		serverName: host,
		localName:  "localhost",
	}

	return c, nil
}

func (c *Client) Close() error {
	return c.Text.Close()
}

func (c *Client) cmd(expectCode int, format string, args ...interface{}) (code int, msg string, err error) {

	id, err := c.Text.Cmd(format, args...)

	if err != nil {
		return 0, "", err
	}

	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)

	code, msg, err = c.Text.ReadResponse(expectCode)

	return code, msg, err

}

func (c *Client) ehlo() error {
	_, msg, err := c.cmd(250, "EHLO %s", c.localName)
	if err != nil {
		return err
	}

	extList := strings.Split(msg, "\n")

	// var ext map[string]string
	// must use make here.
	ext := make(map[string]string)

	if len(extList) > 1 {
		extList = extList[1:]

		for _, line := range extList {
			args := strings.SplitN(line, " ", 2)

			if len(args) > 1 {
				ext[args[0]] = args[1]
			} else {
				ext[args[0]] = ""
			}
		}
	}

	if mechs, ok := ext["AUTH"]; ok {
		c.auth = strings.Split(mechs, " ")
	}
	c.ext = ext
	for k, v := range c.ext {
		fmt.Println(k, v)
	}
	return err

}

func (c *Client) helo() error {
	c.ext = nil
	_, _, err := c.cmd(250, "HELO %s", c.localName)
	return err
}

func (c *Client) Hello(localName string) error {

	err := validateLine(localName)
	if err != nil {
		return err
	}

	if c.didHello {
		return errors.New("smtp: Hello called after other methods")
	}

	c.localName = localName

	return c.hello()
}

func (c *Client) hello() error {

	if !c.didHello {
		c.didHello = true
		err := c.ehlo()
		if err != nil {
			c.helloError = c.helo()
		}
	}
	return c.helloError
}

func (c *Client) Quit() error {

	if err := c.hello(); err != nil {
		return err
	}
	_, _, err := c.cmd(221, "QUIT")
	if err != nil {
		return err
	}
	return c.Close()
}

func (c *Client) Mail(from string) error {
	if err := validateLine(from); err != nil {
		return err
	}

	if err := c.hello(); err != nil {
		return err
	}

	cmdStr := "MAIL FROM:<%s>"

	if c.ext != nil {
		if _, ok := c.ext["8BITMIME"]; ok {
			cmdStr += " BODY=8BITMIME"
		}
	}
	_, _, err := c.cmd(250, cmdStr, from)
	return err
}

func validateLine(line string) error {
	if strings.ContainsAny(line, "\r\n") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func clientFunc(host string) error {

	c, err := Dial(host)
	if err != nil {
		return err
	}

	for {

		var line string
		fmt.Scanf("%s\n", &line)

		cmd, args := parseLine(line)

		switch cmd {
		case "HELO", "EHLO", "HELLO":
			if err := c.Hello(args); err != nil {
				log.Println(err)
			}
		case "MAIL":
			if err := c.Mail(args); err != nil {
				log.Print(err)
				return err
			}

		case "QUIT":
			if err := c.Quit(); err != nil {
				log.Println(err)
				return err
			}
			return nil
		}
	}

}

// server

type Server struct {
	Addr     string
	Hostname string
}

func ListenAndServe(addr string) error {
	s := &Server{Addr: addr}
	return s.ListenAndServe()
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":25"
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(l)

}

func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {

		conn, err := l.Accept()
		if err != nil {
			return err
		}

		ss := s.newSession(conn)

		go ss.Serve()

	}

}

func (s *Server) hostname() string {
	if s.Hostname != "" {
		return s.Hostname
	}
	out, err := exec.Command("hostname").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// The server doesn't need textproto.Conn, but need the textproto.Reader and textproto.Writer
type Session struct {
	s     *Server
	conn  net.Conn
	textr *textproto.Reader
	textw *textproto.Writer

	helloType string
	helloHost string
}

func (s *Server) newSession(conn net.Conn) *Session {
	ss := &Session{
		s:     s,
		conn:  conn,
		textr: textproto.NewReader(bufio.NewReader(conn)),
		textw: textproto.NewWriter(bufio.NewWriter(conn)),
	}
	return ss
}

func (ss *Session) Serve() {
	defer ss.conn.Close()

	// func (w *Writer) PrintfLine(format string, args ...interface{}) error
	ss.textw.PrintfLine("220 %s", ss.s.hostname())
	for {
		line, err := ss.textr.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		cmd, args := parseLine(line)

		switch cmd {

		case "HELO", "EHLO":
			ss.handleHELO(cmd, args)
		case "QUIT":
			ss.handleQUIT()
			return
		default:
			ss.handleDefault()
		}

	}

}

func parseLine(line string) (string, string) {
	if idx := strings.Index(line, " "); idx != -1 {
		cmd := strings.ToUpper(line[:idx])
		args := strings.TrimSpace(line[idx+1:])
		return cmd, args
	}
	return strings.ToUpper(line), ""
}

func (ss *Session) handleHELO(cmd, args string) {

	ss.helloType = cmd
	ss.helloHost = args

	ext := "250-PIPELINING\n" +
		"250-SIZE 10240000\n" +
		"250-ENHANCEDSTATUSCODES\n" +
		"250-8BITMIME\n" +
		"250 DSN"

	err := ss.textw.PrintfLine("250-%s\n%s", ss.s.hostname(), ext)
	if err != nil {
		log.Println(err)
		return
	}

}

func (ss *Session) handleQUIT() {
	ss.textw.PrintfLine("221 Bye")
}

func (ss *Session) handleDefault() {
	ss.textw.PrintfLine("502 5.5.2 Error: command not recognized")
}
