package main

import (
	"encoding/json"
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

func (p Person) String() string {
	s := p.Name.Personal + " " + p.Name.Family

	for _, v := range p.Email {
		s += "\n" + v.Kind + ":" + v.Address
	}

	return s

}

func main() {

	filename := flag.String("f", "", "filename")
	flag.Parse()

	var person Person

	err := loadJSON(*filename, &person)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(person)

}

func loadJSON(filename string, v interface{}) error {

	f, err := os.Open(filename)

	defer f.Close()

	if err != nil {
		return err
	}

	dec := json.NewDecoder(f)

	err = dec.Decode(v)

	return err

}
