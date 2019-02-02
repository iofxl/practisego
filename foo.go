package main

import "fmt"

func main() {

	N := 0x10
	for i := 0; i < N; i++ {
		for j := 0; j <= i; j++ {
			fmt.Printf("%x+%x=%x\t", i, j, i+j)
		}
		fmt.Println()
	}

}
