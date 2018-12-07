package main

import "fmt"
import "math"

func IsPrime(n int, ch chan int) {

	sqrtn := int(math.Sqrt(float64(n)))

	if n%2 == 0 {
		ch <- 0
		return
	}

	for i := 3; i <= sqrtn; i += 2 {
		if n%i == 0 {
			ch <- 0
			return
		}
	}
	ch <- n
}

func main() {

	ch := make(chan int, 10000)
	sum := 5

	for i := 5; i < 2000000; i += 2 {
		go IsPrime(i, ch)
		sum += <-ch
	}

	fmt.Println(sum)

}
