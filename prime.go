// A concurrent prime sieve
// 1. 会生成Gennerate 1个goroutine, 10个 Filter的goroutine.
// 2. 每个loop都取出上个loop后,ch里的值,第一个loop是取出Generate生成的第一个值,就是2
// 3. 有一个要点是,每个loop只print一次,但Filter会一直会一直活着,里面是个永久loop,
// 所以比如,generate生成11这个值后,会依次传給每一层loop,直到某层的loop还没有打印过.
package main

import "fmt"
import "time"

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		fmt.Println("generate: ", i)
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Receive value from 'in'.
		fmt.Printf("my prime is %v, sleep 2s\n", prime)
		time.Sleep(2 * time.Second)
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
		}
	}
}

// The prime sieve: Daisy-chain Filter processes.
func main() {
	ch := make(chan int) // Create a new channel.
	go Generate(ch)      // Launch Generate goroutine.
	for i := 0; i < 10; i++ {
		prime := <-ch
		fmt.Printf("loop %v: %v\n", i, prime)
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
}
