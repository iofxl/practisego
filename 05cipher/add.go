package main

import (
	"fmt"
)

func main() {

	b := make([]byte, 128)

	for i, N := 0, 1000000; i < N; i++ {

		increment(b)
		increment(b)
	}
}

func increment(b []byte) {
	for i := range b {
		b[i]++
		fmt.Println(b)
		if b[i] != 0 {
			return
		}
	}
}
