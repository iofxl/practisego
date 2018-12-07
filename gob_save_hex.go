package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
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

	person := Person{
		Name: Name{"Adel", "Leda"},
		Email: []Email{
			Email{"work", "adel@work.com"},
			Email{"home", "adel@home.com"},
		},
	}

	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)

	err := enc.Encode(&person)
	if err != nil {
		log.Fatal(err)
	}

	bufstring := hex.EncodeToString(buf.Bytes())

	fmt.Println(buf.Bytes())

	fmt.Println(bufstring)

	newb, err := hex.DecodeString(bufstring)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newb)

}
