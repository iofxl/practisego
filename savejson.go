package main

import (
	"encoding/json"
	"flag"
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

	err := saveJSON(*filename, person)

	if err != nil {
		log.Fatal(err)
	}

}

func saveJSON(filename string, v interface{}) error {

	f, err := os.Create(filename)

	defer f.Close()

	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)

	enc.SetIndent("", "\t")

	err = enc.Encode(v)

	return err

}
