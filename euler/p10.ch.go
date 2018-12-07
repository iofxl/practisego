package main

import (
	"fmt"
	"math"
)

//func Sqrt(x float64) float64
func isprime(n int) bool {
	i := math.Sqrt(float64(n))
	for j := 2; float64(j) <= i; j++ {
		if n%j == 0 {
			return false
		}
	}
	return true
}

func main() {

	sum := 0
	for i := 2; i < 2000000000; i++ {
		if isprime(i) {
			fmt.Println(i)
			sum += i
		}

	}

	fmt.Println(sum)

}
