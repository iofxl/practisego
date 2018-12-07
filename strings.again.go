package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {

	str := "Hello World!大哥,洗脚不.\n"
	s := strings.Repeat(str, 2)
	// IO: 6
	// r.I: R, S, WT, RA, ByteScanner, RuneScanner
	sr := strings.NewReader(s)

	b := make([]byte, 3)
	n, _ := sr.Read(b)
	print(b[:n])
	i64, _ := sr.Seek(1, 1)
	// cover auto
	n, _ = sr.Read(b)
	print(i64, b[:n])
	sr.WriteTo(os.Stdout)
	sr.Seek(0, 0)
	sr.ReadAt(b, 3)
	print(b[:n])
	sr.Seek(0, 0)
	by, _ := sr.ReadByte()
	b[0] = by
	print(b[:1])
	i64, _ = sr.Seek(0, 1)
	print(i64)

	by, _ = sr.ReadByte()
	b[0] = by
	print(b[:1])
	i64, _ = sr.Seek(0, 1)
	print(i64)

	_ = sr.UnreadByte()
	i64, _ = sr.Seek(0, 1)
	print(i64)
	_ = sr.UnreadByte()

	i64, _ = sr.Seek(0, 1)
	print(i64)

	i64, _ = sr.Seek(12, 1)
	// "Hello World!大哥,洗脚不."
	print(i64)
	ch, size, _ := sr.ReadRune()
	utf8.EncodeRune(b, ch)
	fmt.Println(string(b), size)
	ch, size, _ = sr.ReadRune()
	utf8.EncodeRune(b, ch)
	fmt.Println(string(b), size)
	_ = sr.UnreadRune()
	i64, _ = sr.Seek(0, 1)
	print(i64)

	l := sr.Len()
	sz := sr.Size()
	print(l, sz)

	s1 := "another string"

	sr.Reset(s1)
	sz = sr.Size()
	print(l, sz)
	sr.WriteTo(os.Stdout)

	rep := strings.NewReplacer("He", "Heeeeeeeeee", "大哥", "老板")

	rep1 := rep.Replace(s)
	n, _ = rep.WriteString(os.Stdout, s)
	print(rep1, n)

	sb := new(strings.Builder)

	sb.Write([]byte("HoooooooooooooooooooooooDo"))
	sb.WriteByte('\t')
	sb.WriteString(s)
	sb.WriteRune(rune('大'))
	print("test sb:", sb.String(), sb.Len())

}

func print(e ...interface{}) {

	for _, v := range e {
		switch v.(type) {
		case int:
			fmt.Printf("%v\n", v)
		case int64:
			fmt.Printf("%v\n", v)
		case string:
			fmt.Printf("%s\n", v)
		case []byte:
			fmt.Printf("%s\n", v)
		case byte:
			fmt.Printf("%x\n", v)
		case rune:
			fmt.Printf("%v\n", v)

		}
	}
}
