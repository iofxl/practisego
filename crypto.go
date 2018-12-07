package main

import (
	"crypto"
	"crypto/sha256"
	"fmt"
	"log"
)

func main() {

	h1 := sha256.New()

	sb := []byte("Hello World!")

	_, err := h1.Write(sb)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x\n", h1.Sum(nil))

	h2 := crypto.SHA256.New()

	_, err = h2.Write(sb)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x\n", h2.Sum(nil))

}
