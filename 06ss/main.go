package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	isserver := flag.Bool("s", false, "server mode")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usage %s address\n", os.Args[0])
	}

	address := flag.Arg(0)

	if *isserver {

		err := ListenAndServeS("tcp4", address)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := ListenAndServe("tcp4", address)

		if err != nil {
			log.Fatal(err)
		}

	}

}
