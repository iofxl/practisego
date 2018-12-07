package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usage %s address\n", os.Args[0])
	}

	address := flag.Arg(0)

	err := ListenAndServe("tcp4", address)

	if err != nil {
		log.Fatal(err)
	}

}
