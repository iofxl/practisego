package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Person struct {
	Name  string `json:"Name"`
	Age   int    `json:"Age"`
	Email string `json:"Email"`
}

func main() {

	var p = Person{"jay", 94, "joo@3z.com"}

	msg, err := json.Marshal(p)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(msg))

	res, err := http.Post("http://localhost/post", "application/json", bytes.NewBuffer(msg))

	if err != nil {
		log.Fatal(err)
	}

	var result map[string]string

	json.NewDecoder(res.Body).Decode(&result)

	var newp Person

	err = json.Unmarshal([]byte(result["data"]), &newp)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newp)

}
