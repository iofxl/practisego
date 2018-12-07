package main

import (
	"fmt"
	"path"
	"strings"
)

func main() {

	a := "  "

	b := " dir2/aa"

	c := "dir3"

	fmt.Println(Join(a, c, b))

}

func Join(elem ...string) string {

	for i, e := range elem {

		fmt.Println(elem)

		fmt.Println(i, e)

		if e != "" {

			return path.Clean(strings.Join(elem[i:], "/"))

		}

	}

	return ""

}
