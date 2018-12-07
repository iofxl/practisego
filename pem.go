package main

import (
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

/*
type Block
func Decode(data []byte) (p *Block, rest []byte)
func Encode(out io.Writer, b *Block) error
func EncodeToMemory(b *Block) []byte
*/

func main() {

	b := pem.Block{"a type", map[string]string{"header": "header 1"}, []byte("Hello World!")}

	err := pem.Encode(os.Stdout, &b)
	// output:
	/*
		-----BEGIN a type-----
		header: header 1

		SGVsbG8gV29ybGQh
		-----END a type-----
	*/

	if err != nil {
		log.Fatal(err)
	}

	sb := pem.EncodeToMemory(&b)

	fmt.Println(string(sb))

	p, rest := pem.Decode(sb)

	fmt.Println(p, string(p.Bytes), "rest:"+string(rest))

}
