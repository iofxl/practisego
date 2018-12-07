package main

import "fmt"
import "unicode/utf8"

func main() {

	// func DecodeLastRune(p []byte) (r rune, size int)
	// func DecodeLastRuneInString(s string) (r rune, size int)
	// func DecodeRune(p []byte) (r rune, size int)
	// func DecodeRuneInString(s string) (r rune, size int)
	b := []byte("hello, 大王!")
	for len(b) >= 0 {

		r, size := utf8.DecodeLastRune(b)
		fmt.Printf("%c %v\n", r, size)

		if len(b) == 0 {
			break
		}

		b = b[:len(b)-size]
	}

	str := "hello, 大王!"
	for len(str) >= 0 {

		r, size := utf8.DecodeLastRuneInString(str)
		fmt.Printf("%c %v\n", r, size)

		if len(str) == 0 {
			break
		}

		str = str[:len(str)-size]
	}

	b = []byte("hello, 大王!")
	for len(b) >= 0 {

		r, size := utf8.DecodeRune(b)
		fmt.Printf("%c %v\n", r, size)

		if len(b) == 0 {
			break
		}

		b = b[size:]
	}

	str = "hello, 大王!"
	for len(str) >= 0 {

		r, size := utf8.DecodeRuneInString(str)
		fmt.Printf("%c %v\n", r, size)

		if len(str) == 0 {
			break
		}

		str = str[size:]
	}

	// func EncodeRune(p []byte, r rune) int
	// if r := "王" { ./utf8.go:68:22: cannot use r (type string) as type rune in argument to utf8.EncodeRune }
	r := '王'
	p := make([]byte, 3)
	n := utf8.EncodeRune(p, r)

	fmt.Printf("%x %v\n", p, n)
	// e78e8b 3

	// FullRune reports whether the bytes in p begin with a full UTF-8 encoding of a rune.
	// An invalid encoding is considered a full Rune since it will convert as a width-1 error rune.
	// func FullRune(p []byte) bool
	// func FullRuneInString(s string) bool
	p = []byte("王")
	fmt.Println(utf8.FullRune(p))
	fmt.Println(utf8.FullRune(p[:2]))
	fmt.Println(utf8.FullRune(p[1:]))
	// true
	// false
	// true    **this one is true, wired.**

	str = "哈喽,大王!"
	fmt.Println(utf8.FullRuneInString(str))
	fmt.Println(utf8.FullRuneInString(str[:2]))
	// true
	// false

	// func RuneCount(p []byte) int
	// func RuneCountInString(s string) (n int)
	p = []byte("哈喽,大王!")
	fmt.Println(len(p))
	fmt.Println(utf8.RuneCount(p))
	// 14
	// 6

	str = "哈喽,大王!"
	fmt.Println(len(str))
	fmt.Println(utf8.RuneCountInString(str))
	// 14
	// 6

	// func RuneLen(r rune) int
	r = '王'
	fmt.Println(utf8.RuneLen(r))
	// 3

	// func RuneStart(b byte) bool
	// Second and subsequent bytes always have the top two bits set to 10.
	p = []byte("a王!")
	for _, b := range p {
		fmt.Println(utf8.RuneStart(b))
	}
	// true
	// true
	// false
	// false
	// true

	// func Valid(p []byte) bool
	// func ValidRune(r rune) bool
	// func ValidString(s string) bool
	p = []byte("a王!")
	invalid := []byte{0xff, 0xfe, 0xfd}
	fmt.Println(utf8.Valid(p))
	fmt.Println(utf8.Valid(invalid))
	// true
	// false

	r = '王'
	invalidr := rune(0xfffffff)
	fmt.Println(utf8.ValidRune(r))
	fmt.Println(utf8.ValidRune(invalidr))
	// true
	// false

	str = "a王!"
	invalids := string([]byte{0xff, 0xfe, 0xfd})
	fmt.Println(utf8.ValidString(str))
	fmt.Println(utf8.ValidString(invalids))
	// true
	// false

}
