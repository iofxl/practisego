package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Message struct {
	Name string
	Body string
	Time int64
}

func main() {

	m := Message{"Adel", "Hello Adel", 123456789123}

	b, err := json.Marshal(m)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	var m2 Message

	err = json.Unmarshal(b, &m2)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(m2)

	b = []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)

	var arbitrary interface{}

	err = json.Unmarshal(b, &arbitrary)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(arbitrary)

}
