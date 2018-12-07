package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/hkdf"
)

func main() {

	secret := make([]byte, 16)

	_, err := io.ReadFull(rand.Reader, secret)

	if err != nil {
		log.Fatal(err)
	}

	salt := make([]byte, 16)

	_, err = io.ReadFull(rand.Reader, salt)

	if err != nil {
		log.Fatal(err)
	}

	info := make([]byte, 5)

	_, err = io.ReadFull(rand.Reader, info)

	if err != nil {
		log.Fatal(err)
	}

	hf := hkdf.New(sha1.New, secret, salt, info)

	key := make([]byte, 16)

	_, err = io.ReadFull(hf, key)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(key))

}
