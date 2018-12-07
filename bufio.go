package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//	"path/filepath"
)

func main() {
	f, err := os.Open("bufio.go")
	defer f.Close()
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		s := scanner.Text()
		s = strings.Replace(s, "f", "FFF", 2)
		fmt.Println(s)
	}

}
