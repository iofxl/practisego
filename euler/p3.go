package main

import (
	"fmt"
	"math"
)

const a int = 600851475143

//const a int = 968

//const a int = 29

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

	sq := math.Sqrt(float64(a))

	result := 0
	for i := 2; float64(i) <= sq; i++ {
		if a%i == 0 {
			if isprime(i) && i > result {
				result = i
			}
			if isprime(a/i) && a/i > result {
				result = a / i
			}

		}
	}

	fmt.Println(result)

}
