package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os/exec"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: func(cmd *cobra.Command, args []string) {
		clientFunc()
	},
}

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {

		if err := ListenAndServe("tcp", addr); err != nil {
			log.Fatal(err)
		}
	},
}

var addr string
var server string

// func (c *Command) AddCommand(cmds ...*Command)
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(serverCmd, clientCmd)
	serverCmd.PersistentFlags().StringVarP(&addr, "listen", "l", "", "listen address")
	clientCmd.Flags().StringVarP(&server, "server", "s", "", "server address")
}

func initConfig() {
	fmt.Println("initConfig")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	Execute()
}

func clientFunc() {

	c, err := Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}

	for {

		var cmd string

		fmt.Scanf("%s\n", &cmd)
		switch cmd {
		case "hello":
			if err := c.hello(); err != nil {
				log.Print(err)
			}

		case "quit":
			if err := c.Quit(); err != nil {
				log.Print(err)
			}
			return
		}

	}

}

type Client struct {
	Text      *textproto.Conn
	localName string
}

func NewClient(conn net.Conn) (*Client, error) {
	text := textproto.NewConn(conn)
	_, _, err := text.ReadResponse(220)
	if err != nil {
		text.Close()
		return nil, err
	}
	c := &Client{Text: text, localName: "localhost"}
	return c, nil

}

func Dial(network, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return NewClient(conn)
}

func (c *Client) cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := c.Text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}

	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)

	code, msg, err := c.Text.ReadCodeLine(expectCode)

	return code, msg, err

}

func (c *Client) hello() error {
	code, msg, err := c.cmd(250, "HELO %s", c.localName)
	fmt.Println(code, msg, err)
	return err
}

func (c *Client) Quit() error {
	if err := c.hello(); err != nil {
		return err
	}
	_, _, err := c.cmd(221, "QUIT")
	if err != nil {
		return err
	}

	return c.Text.Close()
}

type Server struct {
	Network  string
	Addr     string
	Hostname string
	option   string
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":12345"
	}
	l, err := net.Listen(s.Network, addr)
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

func (s *Server) newSession(conn net.Conn) *Session {
	return &Session{
		s:     s,
		conn:  conn,
		textr: textproto.NewReader(bufio.NewReader(conn)),
		textw: textproto.NewWriter(bufio.NewWriter(conn)),
	}
}

func ListenAndServe(network, addr string) error {

	server := &Server{Network: network, Addr: addr, option: ""}

	return server.ListenAndServe()

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
}

func (ss *Session) Serve() {

	defer ss.conn.Close()

	/*
	   // PrintfLine writes the formatted output followed by \r\n.
	       29  func (w *Writer) PrintfLine(format string, args ...interface{}) error {
	       30  	w.closeDot()
	       31  	fmt.Fprintf(w.W, format, args...)
	       32  	w.W.Write(crnl)
	       33  	return w.W.Flush()
	       34  }
	*/
	ss.textw.PrintfLine("220 %s TEXTP textproto", ss.s.hostname())

	for {
		line, err := ss.textr.ReadLine()
		if err != nil {
			log.Print(err)
			return
		}

		var cmd string
		var arg string

		if i := strings.Index(line, " "); i != -1 {
			cmd = strings.ToUpper(line[:i])
			arg = strings.TrimRightFunc(line[i+1:], unicode.IsSpace)
		} else {
			cmd = strings.ToUpper(line)
			arg = ""
		}

		switch cmd {
		case "HELO", "EHLO":
			ss.handleHELO(cmd, arg)
		case "QUIT":
			// s.sendlinef("221 2.0.0 Bye") ; return
			ss.textw.PrintfLine("221 2.0.0 Bye")
			return
		default:
			// s.sendlinef("502 5.5.2 Error: command not recognized")
			ss.textw.PrintfLine("502 5.5.2 Error: command not recognized")
		}

	}

}

// fmt.Fprintf(s.bw, "250-%s\r\n", s.srv.hostname())
func (ss *Session) handleHELO(greeting, host string) {
	fmt.Println("From hello:", greeting, host)
	ss.textw.PrintfLine("250  %s", ss.s.hostname())
}
