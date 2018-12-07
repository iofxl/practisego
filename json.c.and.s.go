package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

type Person struct {
	Name  Name    `json: "name"`
	Email []Email `json: "email"`
}

type Name struct {
	Family   string `json: "family"`
	Personal string `json: "personal"`
}

type Email struct {
	Kind    string `json: "kind"`
	Address string `json: "address:`
}

func (p Person) String() string {

	s := p.Name.Personal + " " + p.Name.Family

	for _, v := range p.Email {
		s += "\n" + v.Kind + ":" + v.Address
	}

	return s
}

func main() {

	client := flag.Bool("c", false, "client")
	flag.Parse()

	address := "127.0.0.1:12345"

	if *client {

		person := Person{
			Name: Name{Family: "Zhang", Personal: "san"},
			Email: []Email{{"work", "zhangsan@work.com"},
				{"home", "zhangsan@home.com"},
			},
		}

		jsonclient(address, person)

		os.Exit(0)
	}

	jsonserver(address)

}

func jsonclient(address string, p Person) error {

	conn, err := net.Dial("tcp4", address)

	defer conn.Close()

	if err != nil {
		return err
	}

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	for i := 0; i < 10; i++ {

		err := enc.Encode(p)

		if err != nil {
			return err
		}

		var p1 Person

		err = dec.Decode(&p1)

		if err != nil {
			return err
		}

		fmt.Println(p1)

	}

	return nil

}

func jsonserver(address string) error {

	l, err := net.Listen("tcp4", address)

	defer l.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {

		conn, err := l.Accept()

		defer conn.Close()

		if err != nil {
			log.Fatal(err)
		}

		enc := json.NewEncoder(conn)
		dec := json.NewDecoder(conn)

		//for i, N := 0, 10; i < N; i++ {
		// 不能做成无限循环的方式
		//for i := 0; ; i++ {

		i := 0
		for dec.More() {

			i++
			fmt.Println(i)

			var p Person

			err = dec.Decode(&p)

			if err != nil {
				continue
			}

			fmt.Println(p)

			err = enc.Encode(p)
			if err != nil {
				continue
			}

		}

		conn.Close()

	}

}
