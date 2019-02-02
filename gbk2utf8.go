package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage: os.Args[0] <gbk_file>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	r := NewGBK2UTF8Reader(f)

	b, err := ioutil.ReadAll(r)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(b))

}

func NewGBK2UTF8Reader(r io.Reader) io.Reader {
	dec := simplifiedchinese.GBK.NewDecoder()
	return dec.Reader(r)
}
