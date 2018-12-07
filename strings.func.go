package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {

	str := "Hello World!大哥,洗脚不."
	s := strings.Repeat(str, 2)
	fmt.Println(len(str)) // 29

	i1 := strings.IndexByte(s, 'o')       // 4
	i2 := strings.Index(s, "不.")          // 25
	i3 := strings.IndexRune(s, rune('哥')) // 15
	i4 := strings.IndexAny(s, "abcdefg")  // 1
	f := func(r rune) bool {
		return unicode.Is(unicode.Han, r)
	}
	i5 := strings.IndexFunc(s, f) // 12

	fmt.Println(i1, i2, i3, i4, i5)

	j1 := strings.LastIndexByte(s, 'o')      // 36
	j2 := strings.LastIndex(s, "不.")         // 54
	j3 := strings.LastIndexAny(s, "abcdefg") // 39 (match is 'd')
	j4 := strings.LastIndexFunc(s, f)        // 54

	fmt.Println(j1, j2, j3, j4)

	k1 := strings.Count(s, "o")              // 4
	b1 := strings.Contains(s, "foo")         // false
	b2 := strings.ContainsRune(s, rune('哥')) // true
	b3 := strings.ContainsAny(s, "abcdefg")  // true
	b4 := strings.HasPrefix(s, "Hoo")        // false
	b5 := strings.HasSuffix(s, "不.")         // true

	fmt.Println(k1, b1, b2, b3, b4, b5)

	ss1 := strings.Split(s, "!")
	ss2 := strings.SplitN(s, "!", 2)
	ss3 := strings.SplitAfter(s, "!")
	ss4 := strings.SplitAfterN(s, "!", 2)
	f = func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsNumber(r))
	}
	ss5 := strings.FieldsFunc(s, f)
	ss6 := strings.Fields(s)

	printrange(ss1, ss2, ss3, ss4, ss5, ss6)

	ss7 := []string{"a", "b", "c"}

	s1 := strings.Join(ss7, ":")
	s2 := strings.Repeat(s1, 7)
	fmt.Println(s1, s2)

	mapping := func(r rune) rune {
		return unicode.ToUpper(r)
	}
	t1 := strings.Map(mapping, s)
	t2 := strings.Replace(s, "洗脚不", "洗不", -1)
	t3 := strings.ToUpper(s)
	t4 := strings.ToLower(s)
	t5 := strings.ToTitle(s)
	t6 := strings.ToUpperSpecial(unicode.TurkishCase, "örnek iş")
	t7 := strings.ToLowerSpecial(unicode.TurkishCase, "örnek iş")
	t8 := strings.ToTitleSpecial(unicode.TurkishCase, "örnek iş")
	t9 := strings.Title(s)

	print(t1, t2, t3, t4, t5, t6, t7, t8, t9)

	f1 := func(r rune) bool {
		return unicode.IsLetter(r)
	}
	f2 := func(r rune) bool {
		return unicode.Is(unicode.Han, r)
	}

	tr1 := strings.TrimLeftFunc(s, f1)
	tr2 := strings.TrimRightFunc(strings.TrimRight(s, "."), f2)
	tr3 := strings.TrimFunc(strings.TrimRight(strings.TrimFunc(s, f1), "."), f2)
	tr4 := strings.TrimLeft(s, "H")
	tr5 := strings.TrimRight(s, "不.")
	tr6 := strings.Trim(strings.Trim(s, "Hello World!"), ",洗脚不.")

	s3 := "     \tHello World!\t\t    "
	tr7 := strings.TrimSpace(s3)
	tr8 := strings.TrimPrefix(s, "Hello")
	// unchanged
	tr9 := strings.TrimSuffix(s, "Hello")

	print(tr1, tr2, tr3, tr4, tr5, tr6, tr7, tr8, tr9)

	c1 := strings.EqualFold("Hello World!", "hello wORLD!") // true
	c2 := strings.Compare("Go", "go")                       // -1

	print(c1, c2)

}

func printrange(e ...[]string) {

	for _, v := range e {

		fmt.Println("Next:")

		for k, vv := range v {
			fmt.Println(k, vv)

		}

	}
}

func print(e ...interface{}) {

	for _, v := range e {

		fmt.Println(v)

	}

}
