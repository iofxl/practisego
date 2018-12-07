package main

import "net/http"
import "log"
import "flag"
import "io"

func main() {

	p := flag.String("p", ":8000", "port")

	flag.Parse()

	foo := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello World!\n")
	}

	http.HandleFunc("/hello", foo)

	err := http.ListenAndServe(*p, nil)

	if err != nil {
		log.Fatal(err)
	}

}
