package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	// func LimitReader(r Reader, n int64) Reader
	r := strings.NewReader("some io.Reader stream to be read\n")
	lr := io.LimitReader(r, 13)

	// output: some io.Reade
	_, err := io.Copy(os.Stdout, lr)

	if err != nil {
		log.Fatal(err)
	}

	p := make([]byte, 2)

	_, err = lr.Read(p)

	if err != nil {
		fmt.Println(err)
	}

	buf := new(bytes.Buffer)
	// func TeeReader(r Reader, w Writer) Reader
	tee := io.TeeReader(r, buf)

	printall := func(r io.Reader) {
		b, err := ioutil.ReadAll(r)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", b)
	}

	printall(tee)
	printall(buf)

	r1 := strings.NewReader("first reader ")
	r2 := strings.NewReader("second reader ")
	r3 := strings.NewReader("third reader\n")

	mr := io.MultiReader(r1, r2, r3)

	_, err = io.Copy(os.Stdout, mr)

	if err != nil {
		log.Fatal(err)
	}

	r.Seek(0, 0)

	buf1, buf2 := new(bytes.Buffer), new(bytes.Buffer)

	mw := io.MultiWriter(buf1, buf2)

	_, err = io.Copy(mw, r)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf1.String())
	fmt.Println(buf2.String())

	pr, pw := io.Pipe()

	go func() {
		fmt.Fprint(pw, "some text to be read\n")
		pw.Close()
	}()

	buf.Reset()

	buf.ReadFrom(pr)
	fmt.Print(buf.String())

	r.Seek(0, 0)

	b := make([]byte, 33)

	_, err = io.ReadAtLeast(r, b, 4)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", b)

	// buffer smaller than minimal read size.
	shortBuf := make([]byte, 3)

	_, err = io.ReadAtLeast(r, shortBuf, 4)

	if err != nil {
		// short buffer
		fmt.Println(err)
	}

	r.Seek(0, 0)

	// minimal read size bigger than io.Reader stream
	longBuf := make([]byte, 64)

	_, err = io.ReadAtLeast(r, longBuf, 64)

	if err != nil {
		// 16:02:10 EOF
		fmt.Println(err)
	}

	fmt.Printf("%s\n", longBuf)

	r.Seek(0, 0)

	sr := io.NewSectionReader(r, 5, 17)

	// io.Reader stream
	_, err = io.Copy(os.Stdout, sr)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Print("\nfoo\n")

	b = b[0:6]

	sr.Seek(0, 0)

	_, err = sr.ReadAt(b, 10)

	// stream
	fmt.Printf("%s\n", b)

	if err != nil {
		log.Fatal(err)
	}

	r.Seek(0, 0)

	rr := rot13Reader{r}

	b = make([]byte, r.Len())

	_, err = rr.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", b)

	buf = bytes.NewBuffer(b)

	rr1 := rot13Reader{buf}

	_, err = io.Copy(os.Stdout, rr1)

	if err != nil {
		log.Fatal(err)
	}

	r = strings.NewReader("Go is a general-purpose language designed with systems programming in mind.")

	b, err = ioutil.ReadAll(r)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", b)

	fi, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	// like `ls -1 | sort`
	for _, f := range fi {
		fmt.Println(f.Name())
	}

	b, err = ioutil.ReadFile("./io2.go")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", b)

	name, err := ioutil.TempDir(".", "temp")

	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(name) // clean up

	f, err := ioutil.TempFile(name, "temp")

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(f.Name()) // clean up

	content := []byte("temporary file's content")

	_, err = f.Write(content)

	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()

	if err != nil {
		log.Fatal(err)
	}

	b, err = ioutil.ReadFile(f.Name())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", b)

	content1 := []byte("temporary file's content by WriteFile")

	err = ioutil.WriteFile(f.Name(), content1, 0666)

	if err != nil {
		log.Fatal(err)
	}

	b, err = ioutil.ReadFile(f.Name())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", b)

}

type rot13Reader struct {
	r io.Reader
}

func (rr rot13Reader) Read(b []byte) (n int, err error) {

	n, err = rr.r.Read(b)

	for i := 0; i < len(b); i++ {

		if b[i] >= 'A' && b[i] < 'N' || b[i] >= 'a' && b[i] < 'n' {
			b[i] += 13
		} else if b[i] >= 'N' && b[i] <= 'Z' || b[i] >= 'n' && b[i] <= 'z' {
			b[i] -= 13
		}
	}
	return
}
