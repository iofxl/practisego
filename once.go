package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once

	onceBody := func() {
		fmt.Println("Only once")
	}

	done := make(chan bool)

	N := 10
	for i := 0; i < N; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}

	for i := 0; i < N; i++ {
		<-done
	}
}
