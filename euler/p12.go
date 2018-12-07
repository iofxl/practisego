package main

import (
	"fmt"
	"math"
)

func factornum(n int) int {

	sum := 0

	if n == 1 {
		return 1
	}

	sqrtn := math.Sqrt(float64(n))
	for i := 1; float64(i) <= sqrtn; i++ {
		if n%i == 0 {
			sum += 2
		}
	}

	return sum

}

func test(in chan int, a int, out chan int) {

	num := 0
	for n := range in {

		num = factornum(n)
		fmt.Println(n, num)
		if num > a {
			out <- n
		}

	}
}

func main() {

	in := make(chan int)
	out := make(chan int)

	for i := 0; i < 1000; i++ {
		go test(in, 500, out)
	}

	go func() {
		sum := 0
		for i := 1; ; i++ {
			sum += i
			in <- sum
		}

	}()

	fmt.Println(<-out)
}
