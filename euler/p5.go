package main

/*

2520 is the smallest number that can be divided by each of the numbers from 1 to 10 without any remainder.

What is the smallest positive number that is evenly divisible by all of the numbers from 1 to 20?
*/

import (
	"fmt"
)

func main() {

	for i := 20; ; i++ {
		if test(i) {
			fmt.Println(i)
			return
		}
	}

}

func test(n int) bool {

	for i := 2; i <= 20; i++ {
		if n%i != 0 {
			return false
		}
	}
	return true
}
