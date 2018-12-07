package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type product interface {
	productup() int
	productdown() int
	productleft() int
	productright() int
	productldiag() int
	productrdiag() int
}

type grid [][]int

func (si grid) productrdiag() int {

	p := 0

	for i, maxi := 0, len(si)-4; i < maxi; i++ {
		for j, maxj := 0, len(si[i])-4; j < maxj; j++ {
			if si[i][j] == 0 || si[i+1][j+1] == 0 || si[i+2][j+2] == 0 || si[i+3][j+3] == 0 {
				continue
			}
			p1 := si[i][j] * si[i+2][j+2] * si[i+3][j+3] * si[i+1][j+1]

			if p1 > p {
				p = p1
			}

		}
	}

	return p
}

func (si grid) productldiag() int {

	p := 0

	for i, mini := len(si)-4, 0; i >= mini; i-- {
		for j, minj := len(si[i])-1, 3; j >= minj; j-- {
			if si[i][j] == 0 || si[i+1][j-1] == 0 || si[i+2][j-2] == 0 || si[i+3][j-3] == 0 {
				continue
			}
			p1 := si[i][j] * si[i+2][j-2] * si[i+3][j-3] * si[i+1][j-1]

			if p1 > p {
				p = p1
			}

		}
	}
	return p
}

func (si grid) productleft() int {

	p := 0

	for i, mini := len(si)-1, 0; i >= mini; i-- {
		for j, minj := len(si[i])-1, 3; j >= minj; j-- {
			if si[i][j] == 0 || si[i][j-1] == 0 || si[i][j-2] == 0 || si[i][j-3] == 0 {
				continue
			}
			p1 := si[i][j] * si[i][j-2] * si[i][j-3] * si[i][j-1]

			if p1 > p {
				p = p1
			}

		}
	}
	return p
}

func (si grid) productright() int {

	p := 0

	for i, mini := len(si)-1, 0; i >= mini; i-- {
		for j, maxj := 0, len(si[i])-4; j < maxj; j++ {
			if si[i][j] == 0 || si[i][j+1] == 0 || si[i][j+2] == 0 || si[i][j+3] == 0 {
				continue
			}
			p1 := si[i][j] * si[i][j+2] * si[i][j+3] * si[i][j+1]

			if p1 > p {
				p = p1
			}

		}
	}
	return p
}

func (si grid) productup() int {

	p := 0

	for i, mini := len(si)-1, 3; i >= mini; i-- {
		for j, maxj := 0, len(si[i]); j < maxj; j++ {
			if si[i][j] == 0 || si[i-3][j] == 0 || si[i-2][j] == 0 || si[i-1][j] == 0 {
				continue
			}
			p1 := si[i][j] * si[i-1][j] * si[i-2][j] * si[i-3][j]

			if p1 > p {
				p = p1
			}

		}
	}
	return p
}

func (si grid) productdown() int {

	p := 0

	for i, mini := len(si)-4, 0; i >= mini; i-- {
		for j, maxj := 0, len(si[i]); j < maxj; j++ {
			if si[i][j] == 0 || si[i+3][j] == 0 || si[i+2][j] == 0 || si[i+1][j] == 0 {
				continue
			}
			p1 := si[i][j] * si[i+1][j] * si[i+2][j] * si[i+3][j]

			if p1 > p {
				p = p1
			}

		}
	}
	return p
}

func maxp(p product) int {

	maxp := p.productleft()

	if p.productright() > maxp {
		maxp = p.productright()
	}

	if p.productup() > maxp {
		maxp = p.productup()
	}
	if p.productdown() > maxp {
		maxp = p.productdown()
	}
	if p.productldiag() > maxp {
		maxp = p.productldiag()
	}
	if p.productrdiag() > maxp {
		maxp = p.productrdiag()
	}

	return maxp
}

// func Open(name string) (*File, error)
func main() {

	f, err := os.Open("p11.txt")

	if err != nil {
		panic(err)
	}

	si := grid{[]int{}}

	s := bufio.NewScanner(f)

	line := 0
	for s.Scan() {

		str := s.Text()

		s1 := bufio.NewScanner(strings.NewReader(str))
		s1.Split(bufio.ScanWords)

		si[line] = []int{}

		for s1.Scan() {
			str := s1.Text()

			i, err := strconv.Atoi(str)

			if err != nil {
				panic(err)
			}

			si[line] = append(si[line], i)
		}
		si = append(si, si[line])
		line++

	}

	fmt.Println(si.productleft())
	fmt.Println(si.productright())
	fmt.Println(si.productup())
	fmt.Println(si.productdown())
	fmt.Println(si.productrdiag())
	fmt.Println(si.productldiag())
	fmt.Println(maxp(si))

}
