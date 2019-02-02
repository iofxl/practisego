package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	fmt.Println(unicode.Is(unicode.L, '的'))
	fmt.Println(unicode.Is(unicode.M, '的'))
	fmt.Println(unicode.Is(unicode.N, '的'))
	fmt.Println(unicode.Is(unicode.P, '的'))
	fmt.Println(unicode.Is(unicode.S, '的'))
	fmt.Println(unicode.Is(unicode.Zs, '的'))
	fmt.Println(unicode.Is(unicode.Han, '的'))

	s := "Hello, 世界"

	for len(s) > 0 {

		r, n := utf8.DecodeLastRuneInString(s)
		fmt.Printf("%c %v\n", r, n)
		s = s[:len(s)-n]
	}

}
