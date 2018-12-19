package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	u, err := url.Parse("http://bing.com/search?q=dotnet")

	if err != nil {
		log.Fatal(err)
	}

	u.Scheme = "https"
	u.Host = "sogo.com"

	q := u.Query()

	q.Set("q", "golang")
	q.Add("foo", "/doo/hoo/")
	u.RawQuery = q.Encode()

	fmt.Println("u is:", u)
	fmt.Println("RawQuery is:", u.RawQuery)

	u, err = url.Parse("ss://foo@doo.com/search?q=any/yna")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(u.String())
	fmt.Println(u.EscapedPath())

}
