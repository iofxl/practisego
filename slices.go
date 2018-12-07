package main

import (
	"fmt"
)

func main() {

	a := []int{1, 2, 3, 4, 5}
	fmt.Println("a: ", a)
	a = append(a, 88)
	fmt.Println("after append 88: ", a)
	a = append(a, a...)
	fmt.Println("after append a...: ", a)
	copy(a, a[1:])
	fmt.Println("after copy a[1:]: ", a)

	//实现插入

	i := 3

	//这个是不安全的,如果len(a) == cap(a) 就要panic了.
	a = a[:len(a)+1]
	copy(a[i+1:], a[i:])
	a[i] = 33
	fmt.Println("after insert a[3]: ", a)

	a = append(a)
	fmt.Println(a)

	fmt.Println(&a[0])
	a = append(a, a...)
	fmt.Println(&a[0])

	a = append([]int(nil), a...)
	fmt.Println(a)

	// cut
	a = []int{1, 2, 3, 4, 5}
	fmt.Println("cut:", a)
	a = append(a[:1], a[1+3:]...)
	fmt.Println("cut:", a)

	// delete
	a = []int{1, 2, 3, 4, 5}
	fmt.Println("delete: ", a)
	a = append(a[:0], a[0+1:]...)
	fmt.Println("delete: ", a)

	a = []int{1, 2, 3, 4, 5}
	fmt.Println(a[1:5])

}
