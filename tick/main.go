package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	n := 4.5

	var stack []chan struct{}

	t := time.NewTicker(2 * time.Second)

	for range t.C {

		fmt.Println("Goroutine before: ", runtime.NumGoroutine())

		if n < 5 {
			n += 1
			stack = append(stack, add())
			fmt.Println("Goroutine after append: ", runtime.NumGoroutine())
		} else {
			n -= 1
			close(stack[len(stack)-1])
			stack = stack[:len(stack)-1]
			fmt.Println("Goroutine after drop: ", runtime.NumGoroutine())
		}
		fmt.Println("Goroutine after: ", runtime.NumGoroutine())
	}

}

type tickers struct {
	ch   []chan struct{}
	idle float64
}

func add() chan struct{} {
	stop := make(chan struct{})
	go func() {
		go tick(stop)
	}()
	return stop
}

func tick(stop chan struct{}) {
	go func() {
		t := time.NewTicker(time.Millisecond)
		for range t.C {
			select {
			case <-stop:
				return
			default:
			}
		}
	}()
}
