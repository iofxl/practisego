package main

import (
	"fmt"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type StructA struct {
	A string `yaml:"a"`
}

type StructB struct {
	StructA `yaml:",inline"`
	B       string `yaml:"b"`
}

func main() {

	data := `
a: a string from StructA
b: a string from StructB
`

	var b StructB

	err := yaml.Unmarshal([]byte(data), &b)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b.A, "\n", b.B)

}
