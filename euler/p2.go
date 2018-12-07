package main

// euler: p2

import (
	"fmt"
)

func main() {

	sum := 0
	for a, b := 1, 2; a <= 4000000; {
		if a%2 == 0 {
			sum += a
		}
		a, b = b, a+b
		fmt.Println(a, b, sum)
	}

	fmt.Println("sum :", sum)

}
