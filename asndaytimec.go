package main

import (
	"encoding/asn1"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)

func main() {

	address := "127.0.0.1:12345"

	conn, err := net.Dial("tcp4", address)

	if err != nil {
		log.Fatal(err)
	}

	mdata, err := ioutil.ReadAll(conn)

	if err != nil {
		log.Fatal(err)
	}

	var daytime time.Time

	_, err = asn1.Unmarshal(mdata, &daytime)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(mdata, daytime)

}
