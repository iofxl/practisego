package main

import (
	"fmt"
	"math/bits"
)

func main() {
	b := []byte{41, 41, 41, 41, 41, 41, 41, 41}
	for _, v := range b {
		fmt.Printf("%8b,%8b,%8b\n", v, bits.Reverse8(v<<4), v>>4)
	}
}
