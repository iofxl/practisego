package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
)

func main() {

	b := make([]byte, 10)

	rand.Reader.Read(b)
	fmt.Printf("%v\n", b)

	b = make([]byte, 10)
	rand.Read(b)
	fmt.Printf("%v\n", b)

	sr := strings.NewReader("Hello World!")

	bi := big.NewInt(10)

	n, err := rand.Int(sr, bi)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(n)

	p, err := rand.Prime(sr, 5)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p)
}
