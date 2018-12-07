package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type P struct {
	X, Y, Z int
	Name    string
}

func decgob(w http.ResponseWriter, r *http.Request) {

	var p P

	dec := gob.NewDecoder(r.Body)

	dec.Decode(&p)

	fmt.Println(p)

}

func main() {

	var addr string

	flag.StringVar(&addr, "a", ":12345", "addr")
	flag.Parse()

	http.HandleFunc("/post/", decgob)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal(err)
	}

}
