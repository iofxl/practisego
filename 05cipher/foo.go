package main

import "fmt"

func main() {

	var foo []byte

	for {

		fmt.Scanf("%v\n", &foo)

		fmt.Println(foo, string(foo))
	}
}
