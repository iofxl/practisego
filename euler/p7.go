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

	i := 2
	for count := 0; count < 10001; i++ {
		if isprime(i) {
			count++
		}

	}

	fmt.Println(i - 1)

}
