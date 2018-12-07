// https://github.com/protocolbuffers/protobuf/tree/master/examples
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	pb "./pb"
	"github.com/golang/protobuf/proto"
)

// 读出整个文件，生成东西，添加东西，写转文件
func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s foo.book\n", os.Args[0])
	}

	fname := os.Args[1]

	in, err := ioutil.ReadFile(fname)

	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("%s is not exit\n", fname)
		} else {
			log.Fatalf("open %s error\n", fname)
		}
	}

	book := new(pb.AddressBook)

	err = proto.Unmarshal(in, book)

	if err != nil {
		log.Fatal(err)
	}

	people, err := promptForAdd(os.Stdin)

	if err != nil {
		log.Fatal(err)
	}

	// book是一个pb.AddressBook的指针，book的field是People, 而不中Person, Person是People结构的类型，或者曰是结构名
	// field名相当于变量名，不过变量赋值使"=", field赋值使":"
	book.People = append(book.People, people)

	out, err := proto.Marshal(book)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fname, out, 0644)

	if err != nil {
		log.Fatal(err)
	}

}

func promptForAdd(r io.Reader) (*pb.Person, error) {
	fmt.Println("Enter name:")

	rd := bufio.NewReader(r)

	p := new(pb.Person)

	_, _ = rd.ReadString('\n')

	p = &pb.Person{
		Name:  "zhangsan",
		Id:    1,
		Email: "zhangsan@jp.com",

		Phones: []*pb.Person_PhoneNumber{
			{Number: "110", Type: pb.Person_HOME},
		},
	}

	return p, nil

}
