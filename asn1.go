package main

import (
	"encoding/asn1"
	"fmt"
	"log"
)

func main() {

	a := 13
	var b int

	buf, err := asn1.Marshal(a)

	if err != nil {
		log.Fatal(err)
	}

	_, err = asn1.Unmarshal(buf, &b)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b)
}
