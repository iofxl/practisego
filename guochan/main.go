package main

import (
	"crypto/rand"
	"flag"
	"io"
	"log"
	"os"
)

func main() {

	isserver := flag.Bool("s", false, "Server mode")
	address := flag.String("l", ":443", "Listen address: IP:Port")
	method := flag.Uint("m", 1, "Method")
	flag.Parse()

	cfg, err := loadConfig("guochan.cfg")

	if err != nil {

		if !*isserver {
			if flag.NArg() != 1 {
				log.Fatalf("Usage: %s server_string\n", os.Args[0])
			}
			s := flag.Arg(0)
			cfg, err := decodeConfig(s)
			if err != nil {
				log.Fatal(err)
			}

			cfg.Address = *address
			err = saveConfig(cfg, "guochan.cfg")
			if err != nil {
				log.Fatal(err)
			}

		} else {

			cfg := &Config{Secret: make([]byte, 16), Salt: make([]byte, 16), Info: make([]byte, 16)}
			io.ReadFull(rand.Reader, cfg.Secret)
			io.ReadFull(rand.Reader, cfg.Salt)
			io.ReadFull(rand.Reader, cfg.Info)
			cfg.M = Method(*method)
			cfg.Address = *address
			cfg.Server = cfg.Address
			if err := saveConfig(cfg, "guochan.cfg"); err != nil {
				log.Fatal("saveConfig error:", err)
			}

		}
	}

	if *isserver {
		ListenAndServeS(cfg)
	} else {
		ListenAndServeC(cfg)
	}

}
