package main

import "fmt"
import "math"

type person struct {
	xing string
	ming string
}

type tewu struct {
	person
	ltk bool
}

type square struct {
	side float64
}

type circle struct {
	radius float64
}

func (p person) psay() {
	fmt.Println(p.ming, `説,"天王盖地虎!"`)
}

func (t tewu) tsay() {
	fmt.Println(t.person.ming, `説,"滚犊子!"`)
	fmt.Println("has a ltk:", t.ltk)
}

func (s square) area() float64 {
	return s.side * s.side
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

type shape interface {
	area() float64
}

func info(s shape) {
	fmt.Println(s.area())
}

func main() {

	p := person{"张", "三"}
	t := tewu{person{"李", "四"}, true}
	fmt.Println(p.xing, p.ming)
	fmt.Println(t.person.xing, t.person.ming, t.ltk)

	p.psay()
	t.tsay()
	t.psay()

	s := square{2}
	c := circle{1}

	info(s)
	info(c)

}
