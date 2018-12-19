package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	client := http.Client{}

	req, err := http.NewRequest("GET", "http://localhost/get", nil)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	fmt.Println("data:", result["data"])
}
