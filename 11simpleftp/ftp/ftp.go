package ftp

import (
	"errors"
	"net"
	"net/textproto"
	"strings"
)

const (
	StatusReady = 220
)

type Client struct {
	Text textproto.Conn
}

func NewClient(conn net.Conn) (*Client, error) {

	text := textproto.NewConn(conn)
	_, _, err := text.ReadResponse(StatusReady)
	if err != nil {
		text.Close()
		return nil, err
	}

	c := &Client{Text: text}
	return c, nil
}

func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

func (c *Client) Close() error {
	return c.Text.Close()
}

func (c *Clinet) cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := c.Text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
	code, msg, err := c.Text.ReadResponse(expectCode)
	return cocde, msg, err

}

// Implement dir cd pwd quit

func (c *Client) Dir() error {
	code, msg, err := c.cmd(

}

func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("ftp: A line must not contain CR or LF")
	}
	return nil
}
