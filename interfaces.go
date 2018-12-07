package main

import (
	"fmt"
	"math"
)

//an interface
type areaer interface {
	area() float64
}

type rect struct {
	width, height float64
}

type circle struct {
	radius float64
}

func (r rect) area() float64 {
	return r.width * r.height
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func measure(a areaer) float64 {
	return a.area()
}

func main() {

	r1 := rect{width: 5, height: 3}
	c1 := circle{radius: 10}

	// type rect(circlt) has area(), so can use it as type areaer,but why ?
	fmt.Println(measure(r1))
	fmt.Println(measure(c1))
}
