package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"net"

	"golang.org/x/crypto/hkdf"
)

// goss2的做法是通过AESCTR函数生成blk, 把blk使ctrStream结构包起,这个ctrStream实现了这个包里的Cipher接口
// 这个Cipher接口的作用,也就是生成加密的工作模式, 所以相当于块的

// 这里要怎样从cfg 转化到这个Method里去呢？

func NewCTR(m Method, iv []byte) cipher.Stream {

	secret := cfg.Secret
	salt := cfg.Salt
	info := cfg.Info

	key := make([]byte, keySizes[cfg.M])
	hkdfrd := hkdf.New(sha1.New, secret, salt, info)
	io.ReadFull(hkdfrd, key)

	blk, _ := aes.NewCipher(key)

	ctr := cipher.NewCTR(blk, iv)

	return ctr
}

type conn struct {
	net.Conn
	m Method
	r *cipher.StreamReader
	w *cipher.StreamWriter
}

func NewConn(m Method, c net.Conn) net.Conn {
	return &conn{Conn: c, m: m}
}

// 犯了个傻逼错误, 没有用指针去实现,结果总是已经initWriter()的conn, 要写的时候还是报c.w是nil
func (c *conn) initReader() error {
	if c.r == nil {
		iv := make([]byte, keySizes[c.m])
		if _, err := io.ReadFull(c.Conn, iv); err != nil {
			return err
		}
		ctr := NewCTR(c.m, iv)
		c.r = &cipher.StreamReader{S: ctr, R: c.Conn}
	}
	return nil
}

func (c *conn) Read(p []byte) (int, error) {
	if err := c.initReader(); err != nil {
		return 0, err
	}
	return c.r.Read(p)
}

func (c *conn) initWriter() error {
	if c.w == nil {
		iv := make([]byte, keySizes[c.m])
		io.ReadFull(rand.Reader, iv)
		ctr := NewCTR(c.m, iv)
		if _, err := c.Conn.Write(iv); err != nil {
			return err
		}
		c.w = &cipher.StreamWriter{S: ctr, W: c.Conn}

	}
	return nil
}

func (c *conn) Write(p []byte) (int, error) {

	if err := c.initWriter(); err != nil {
		return 0, err
	}
	return c.w.Write(p)
}
