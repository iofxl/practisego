package main

import "fmt"

func main() {

	for a := 1; a < 1000; a++ {
		for b := 2; b < 1000; b++ {
			if b <= a {
				continue
			}
			for c := 3; c < 1000; c++ {
				if c <= b {
					continue
				}
				if a*a+b*b == c*c && a+b+c == 1000 {
					fmt.Printf("%d\n", a*b*c)
					return
				}
			}
		}
	}

}
