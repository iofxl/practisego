package main

import (
	"fmt"
	"runtime"
	"sync"
	//"time"
)

func main() {
	// runtime.GOMAXPROCS(1)

	// 1 2 3 4 5 6 7 8 9 10 a b c d e f g h i j k l m n o p q r s t u v w x y z 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
	// 1 2 3 4 5 6 7 8 9 10 11 12 13 a b c d 14 15 16 17 18 19 20 21 22 23 24 25 26 e f g h i j k l m n o p q r s t u v w x y z
	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup

	wg.Add(2)

	fmt.Println("Starting Go Routines")
	go func() {
		defer wg.Done()

		for char := 'a'; char < 'a'+26; char++ {
			fmt.Printf("%c ", char)
		}
	}()

	go func() {
		defer wg.Done()

		//	time.Sleep(2 * time.Second)
		for number := 1; number < 27; number++ {
			fmt.Printf("%d ", number)
		}
	}()

	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")
}
