package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {

	h := sha256.New()

	h.Write([]byte("Hello World!"))

	fmt.Printf("%x\n", h.Sum(nil))

	fmt.Printf("%x\n", (sha256.Sum256([]byte("Hello World!"))))

	fmt.Printf("%x\n", (sha256.Sum224([]byte("Hello World!"))))

	h2 := sha256.New224()

	h2.Write([]byte("Hello World!"))

	b := make([]byte, 1)

	sum := h2.Sum(b)

	fmt.Printf("%x\n", sum)

	fmt.Printf("%x\n", b)

}
