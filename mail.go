package main

import (
	"fmt"
	"log"
	"net/mail"
	"strings"
)

func main() {

	msgstring := `Date: 20180915 15:17:43 +0800
From: Gopher <from@example.com>
To: Another Gopher <to@example.com>
Subject: Gophers at Gophercon

Message Body
`
	r := strings.NewReader(msgstring)

	msg, err := mail.ReadMessage(r)

	if err != nil {
		log.Fatal(err)
	}

	h := msg.Header

	fmt.Println("Date:", h.Get("Date"))
	fmt.Println("From:", h.Get("From"))
	fmt.Println("To:", h.Get("To"))
	fmt.Println("Subject:", h.Get("Subject"))

	t, err := h.Date()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(t.String())

}
