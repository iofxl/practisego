package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type P struct {
	X, Y, Z int
	Name    string
}

func main() {

	var url string

	flag.StringVar(&url, "url", "http://127.0.0.1:12345/post/", "url")
	flag.Parse()

	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)

	enc.Encode(P{7, 8, 9, "foo"})

	//func Post(url string, contentType string, body io.Reader) (resp *Response, err error)
	resp, err := http.Post(url, "gob; charset=utf8", buf)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	fmt.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

}
