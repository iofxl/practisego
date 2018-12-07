package main

import "fmt"
import "math"

func sqart(x float64) float64 {
	z := 1.0
	for math.Abs(z*z-x) > 0.0000001 {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

func main() {

	fmt.Println(sqart(2))
}
