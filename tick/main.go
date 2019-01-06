package main

import (
	"flag"
	"log"
	"os"
	"runtime"
)

func main() {

	var cputime uint
	var memused uint

	flag.UintVar(&cputime, "c", 30, "cpu time: 1 - 80 %")
	flag.UintVar(&memused, "r", 30, "mem used: 1 - 80 %")
	flag.Parse()

	if cputime > 80 {
		log.Fatalln("cpu time must <= 80")
	}

	if memused > 80 {
		log.Fatalln("mem used must <= 80")
	}

	t := &ticker{nil, make(chan float64), cputime, runtime.NumCPU()}
	go t.run()

	r := &mem{nil, make(chan float64), 0, memused, os.Getpagesize()}
	r.run()

}
