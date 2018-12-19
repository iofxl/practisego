package main

import (
	"log"
	"net/smtp"
)

func main() {

	addr := "localhost:1025"
	auth := smtp.CRAMMD5Auth("username", "password")
	from := "x@y.org"
	to := []string{"z1@z.org", "z2@z.org"}
	msg := []byte("From: from@y.org\r\n" +
		"To: z1@z.org, z2@z.org\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := smtp.SendMail(addr, auth, from, to, msg)

	if err != nil {
		log.Fatal(err)
	}

}
