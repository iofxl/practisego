package main

import (
	"bytes"
	"fmt"
)

func main() {

	buf := new(bytes.Buffer)

	buf.Write([]byte("1234567"))

	fmt.Println(buf.Cap(), buf.Len())

	p := make([]byte, 8)

	n, err := buf.Read(p)
	fmt.Println(n, err)
	n, err = buf.Read(p)
	fmt.Println(n, err)
	n, err = buf.Read(p)
	fmt.Println(n, err)

	// when content > len(p)
	// 6 <nil>
	// 1 <nil>
	// 0 EOF

	// when content <= len(p)
	// 7 <nil>
	// 0 EOF
	// 0 EOF

	b := make([]byte, 5)

	buf1 := bytes.NewBuffer(b)

	n, err = buf1.Write([]byte("12346"))
	fmt.Println("Write:", n, err)

}
