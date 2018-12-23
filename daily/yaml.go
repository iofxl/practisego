package main

import (
	"fmt"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type A struct {
	B string
	C int
}

type AS struct {
	P As
}

type As []A

func main() {

	var as = As{{"zhang1", 1}, {"zhang2", 2}, {"zhang3", 3}}
	var aS = AS{As{{"zhang1", 1}, {"zhang2", 2}, {"zhang3", 3}}}

	b, err := yaml.Marshal(as)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
	fmt.Println()

	bS, err := yaml.Marshal(aS)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bS))

}
