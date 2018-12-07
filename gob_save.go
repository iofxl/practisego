package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
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

	person := Person{
		Name: Name{"Adel", "Leda"},
		Email: []Email{
			Email{"work", "adel@work.com"},
			Email{"home", "adel@home.com"},
		},
	}

	err := savegob(*filename, person)

	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(*filename)

	content, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content)

}

func savegob(filename string, v interface{}) error {

	f, err := os.Create(filename)

	defer f.Close()

	if err != nil {
		return err
	}

	enc := gob.NewEncoder(f)

	err = enc.Encode(v)

	return err

}
