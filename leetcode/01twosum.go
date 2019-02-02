package main

import "fmt"

func main() {
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
}

func twoSum(nums []int, target int) []int {

	l := len(nums)

	for i := 0; i < l-1; i++ {
		for j := i + 1; j < l; j++ {

			if nums[j] == target-nums[i] {
				return []int{i, j}
			}
		}

	}
	return []int{-1}

}
