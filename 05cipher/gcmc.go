package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/hkdf"
)

type AEADWriter struct {
	A     cipher.AEAD
	W     io.Writer
	nonce []byte
}

func (w AEADWriter) Write(src []byte) (n int, err error) {
	c := make([]byte, len(src)+w.A.Overhead())
	c = w.A.Seal(nil, w.nonce, src, nil)
	n, err = w.W.Write(c)
	return
}

type AEADReader struct {
	A     cipher.AEAD
	R     io.Reader
	nonce []byte
}

func (r AEADReader) Read(dst []byte) (n int, err error) {
	n, err = r.R.Read(dst)
	r.A.Open(dst[:0], r.nonce, dst[:n], nil)
	fmt.Println("In AEADReader:", n, err)
	if err != nil {
		return n, err
	}
	return n - r.A.Overhead(), nil
}

func main() {

	//salt := make([]byte, 16)
	//rand.Read(salt)

	// secret := make([]byte, 128)
	//rand.Read(secret)

	secret := []byte("secret")
	salt := []byte("salt")
	info := []byte("info")
	hkdfrd := hkdf.New(sha1.New, secret, salt, info)

	key := make([]byte, 16)
	_, err := hkdfrd.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	blk, _ := aes.NewCipher(key)
	aead, _ := cipher.NewGCM(blk)

	conn, err := net.Dial("tcp4", ":12345")

	if err != nil {
		log.Fatal(err)
	}

	handleConn(conn, aead)

}

func handleConn(conn net.Conn, aead cipher.AEAD) {

	defer conn.Close()

	nonce := make([]byte, 12)

	aw := &AEADWriter{A: aead, W: conn, nonce: nonce}
	ar := &AEADReader{A: aead, R: conn, nonce: nonce}

	var src []byte

	resp := make([]byte, 4096)

	for {

		n, _ := fmt.Scanf("%v\n", &src)
		fmt.Println("src:", src, string(src))

		n, err := aw.Write(src)
		fmt.Println("Write length:", n)

		if err != nil {
			log.Println("aw Writer error:", err)
			continue
		}

		n, err = ar.Read(resp)
		fmt.Println("Read length:", n)

		if err != nil {
			if err == io.EOF {
				log.Println("ar Read error", err)
				continue
			}
			log.Println("ar Read error", err)
			return
		}

		fmt.Println("message:", string(resp[:n]))

	}

}
