package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
)

func main() {

	f, err := ioutil.ReadFile("p13.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	// func Split(s, sep []byte) [][]byte
	ff := bytes.Split(f, []byte{'\n'})

	m, result := 0, make([]byte, 0)

	for j := 49; j >= 0; j-- {
		for i := 0; i < 100; i++ {
			// func Atoi(s string) (int, error)
			v, _ := strconv.Atoi(string(ff[i][j]))
			m += v
		}
		// func Itoa(i int) string
		result = append(result, []byte(strconv.Itoa(m%10))...)
		m /= 10

	}

	for m > 0 {
		result = append(result, []byte(strconv.Itoa(m%10))...)
		m /= 10
	}

	for i, n := len(result)-1, len(result)-10; i >= n; i-- {
		fmt.Print(string(result[i]))
	}
}
