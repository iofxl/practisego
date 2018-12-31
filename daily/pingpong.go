package main

import (
	"fmt"
	"time"
)

func main() {
	var Ball int
	table := make(chan int)
	for i := 0; i < 10; i++ {
		go player(table)
	}

	table <- Ball
	time.Sleep(100 * time.Second)
	<-table

}

func player(table chan int) {
	for {
		ball := <-table
		ball++
		time.Sleep(100 * time.Millisecond)
		table <- ball
		fmt.Println(ball)
	}
}
