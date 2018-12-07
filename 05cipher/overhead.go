package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/hkdf"
)

func main() {

	secret, salt, key := make([]byte, aes.BlockSize), make([]byte, aes.BlockSize), make([]byte, aes.BlockSize)

	_, err := io.ReadFull(rand.Reader, secret)

	if err != nil {
		log.Fatal(err)
	}

	_, err = io.ReadFull(rand.Reader, salt)

	if err != nil {
		log.Fatal(err)
	}

	hkdfrd := hkdf.New(sha1.New, secret, salt, []byte("info"))

	io.ReadFull(hkdfrd, key)

	blk, err := aes.NewCipher(key)

	if err != nil {
		log.Fatal(err)
	}

	aead, err := cipher.NewGCM(blk)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NonceSize:", aead.NonceSize())
	fmt.Println("Overhead:", aead.Overhead())
}
