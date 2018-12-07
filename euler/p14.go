package main

import (
	"fmt"
	"sync"
)

func gen() <-chan int {

	ch := make(chan int)

	go func() {
		for i := 1; i <= 1000000; i++ {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func coll(numc <-chan int) chan map[int]int {

	out := make(chan map[int]int)
	var wg sync.WaitGroup

	output := func(numc <-chan int) {
		defer wg.Done()

		for nn := range numc {
			n := nn
			m, chain := make(map[int]int), 1
			for n != 1 {
				if n%2 == 0 {
					n /= 2
					chain++
				} else {
					n = 3*n + 1
					chain++
				}
			}
			m[nn] = chain
			out <- m
		}
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go output(numc)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out

}

func main() {

	numc := gen()

	mc := coll(numc)

	max, maxk := 1, 1

	for m := range mc {
		fmt.Println(m)
		for k, v := range m {
			if v > max {
				max = v
				maxk = k
			}
		}
	}

	fmt.Println(maxk)

}
