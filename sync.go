package main

import (
	"fmt"
	"sync"
)

type SafeInt struct {
	sync.Mutex
	Num int
}

func main() {
	count := SafeInt{}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(i int) {
			count.Lock() // 加锁，防止其它例程修改 count
			count.Num += i
			fmt.Print(count.Num, " ")
			count.Unlock() // 修改完毕，解锁
			done <- true
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
