package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
)

type Person struct {
	Name  Name
	Email []Email
}

type Name struct {
	Family   string
	Personal string
}

type Email struct {
	Kind    string
	Address string
}

func main() {

	filename := flag.String("f", "", "filename")
	flag.Parse()

	var person Person

	err := loadGOB(*filename, &person)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(person)

}

func loadGOB(filename string, v interface{}) error {

	f, err := os.Open(filename)

	defer f.Close()

	if err != nil {
		return err
	}

	dec := gob.NewDecoder(f)

	err = dec.Decode(v)

	return err

}
