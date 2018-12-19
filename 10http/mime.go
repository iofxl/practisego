package main

import (
	"fmt"
	"io"
	"mime"
)

func main() {

	s := "Hello, 世界"

	b := mime.BEncoding.Encode("utf-8", s)
	q := mime.QEncoding.Encode("utf-8", s)

	fmt.Println("BEncoding:", b)
	fmt.Println("BEncoding:", q)

	dec := new(mime.WordDecoder)

	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		default:
		}

		return input, nil
	}

	bout, _ := dec.Decode(b)
	bhout, _ := dec.DecodeHeader(b)

	fmt.Println("Decode from b:", bout)
	fmt.Println("DecodeHeader from b:", bhout)

	qout, _ := dec.Decode(q)
	qhout, _ := dec.DecodeHeader(q)

	fmt.Println("Decode from q:", qout)
	fmt.Println("Decode from q:", qhout)

}
