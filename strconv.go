package main

import (
	"fmt"
	"strconv"
)

func main() {

	v32 := "-354634382"

	if i, err := strconv.ParseInt(v32, 10, 32); err == nil {

		fmt.Printf("%T, %v\n", i, i)
	}

	v64 := "-3546343826724305832"

	if i, err := strconv.ParseInt(v64, 10, 64); err == nil {

		fmt.Printf("%T, %v\n", i, i)
	}

}
