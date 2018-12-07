package main

import (
	"fmt"
	"sort"
)

// Largest palindrome product Problem 4
// A palindromic number reads the same both ways. The largest palindrome made from the product of two 2-digit numbers is 9009 = 91 × 99.
// Find the largest palindrome made from the product of two 3-digit numbers.
func main() {

	// 100 ... 999

	var all []int
	for i := 999; i > 99; i-- {
		for j := 999; j > 99; j-- {
			all = append(all, i*j)
		}
	}

	sort.IntSlice(all).Sort()

	for i := len(all) - 1; i >= 0; i-- {

		if test(split(all[i])) {
			fmt.Println(all[i])
			return
		}
	}

}

func split(n int) []int {

	var b []int

	for n >= 10 {
		b = append(b, n%10)
		n = n / 10
	}
	b = append(b, n)

	return b

}

func test(b []int) bool {
	for left, right := 0, len(b)-1; left < right; left, right = left+1, right-1 {
		if b[left] != b[right] {
			return false
		}
	}
	return true
}
