package main

import (
	"fmt"
	"io/ioutil"

	//	"strings"
	"os"
)

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	//the default is os.TempDir()
	td, err := ioutil.TempDir(os.TempDir(), "testioutil")
	defer os.RemoveAll(td)
	checkerr(err)
	fmt.Println("TempDir is", td)

	f, err := ioutil.TempFile(td, "testfile")
	defer os.Remove(f.Name())
	checkerr(err)
	fmt.Println("TempFile is", f.Name())

	content, err := ioutil.ReadFile("ioutil.go")

	ioutil.WriteFile(f.Name(), content, 0660)
	//check Seek is where now, offset is back to Seek(0,0)
	s, _ := f.Seek(0, 1)
	fmt.Println("f's offset after ioutil.WriteFile:", s)

	buf, err := ioutil.ReadFile(f.Name())
	checkerr(err)
	fmt.Println(string(buf))
	s, _ = f.Seek(0, 1)
	fmt.Println("f's offset after ioutil.ReadFile:", s)

	b := make([]byte, 8)

	n, err := f.Read(b)
	checkerr(err)
	fmt.Println("f.Read:", n, err, string(b))

	s, _ = f.Seek(0, 1)
	fmt.Println("f's offset after f.Read(b)", s)

	//what will happen do ioutil.ReadAll now

	ra, err := ioutil.ReadAll(f)
	checkerr(err)
	fmt.Println("ioutil.ReadAll:", string(ra))

	s, _ = f.Seek(0, 1)
	fmt.Println("f's offset after ioutil.ReadAll", s)

	rd, err := ioutil.ReadDir(td)
	checkerr(err)

	for _, file := range rd {
		fmt.Println(file.Name())
	}

}
