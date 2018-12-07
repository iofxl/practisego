package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	pb "./pb"
	"github.com/golang/protobuf/proto"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("foo.book")
	}

	fname := os.Args[1]

	in, err := ioutil.ReadFile(fname)

	if err != nil {
		log.Fatal(err)
	}

	book := new(pb.AddressBook)

	err = proto.Unmarshal(in, book)

	if err != nil {
		log.Fatal(err)
	}

	listpeople(book)

}

func listpeople(book *pb.AddressBook) {

	for _, v := range book.People {
		fmt.Println(v.GetName(), v.GetId(), v.GetEmail(), v.GetPhones())

		for _, v := range v.Phones {
			fmt.Println(v.GetNumber(), v.GetType())

		}

	}

	for _, v := range book.People {
		fmt.Println(v.String())

		for _, v := range v.Phones {
			fmt.Println(v.String())

		}

	}

}
