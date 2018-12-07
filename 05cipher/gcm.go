package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"

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
	return n - r.A.Overhead(), err
}

func main() {

	secret := make([]byte, 128)
	salt := make([]byte, 16)
	info := []byte("anbababa")
	rand.Read(secret)
	rand.Read(salt)
	hkdfrd := hkdf.New(sha1.New, secret, salt, info)

	key := make([]byte, 16)
	_, err := hkdfrd.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	blk, _ := aes.NewCipher(key)
	aead, _ := cipher.NewGCM(blk)

	buf := new(bytes.Buffer)

	nonce := make([]byte, 12)

	rand.Read(nonce)

	aw := &AEADWriter{A: aead, W: buf, nonce: nonce}
	ar := &AEADReader{A: aead, R: buf, nonce: nonce}

	for {

		var src []byte

		fmt.Scanf("%v\n", &src)
		aw.Write(src)
		fmt.Println("dst:", hex.EncodeToString(aw.W.(*bytes.Buffer).Bytes()))

		dst := make([]byte, 512)
		n, err := ar.Read(dst)

		if err != nil {
			log.Println("ar Read error", err)
			continue
		}
		fmt.Println("src:", string(dst[:n]))

		buf.Reset()

	}

}
