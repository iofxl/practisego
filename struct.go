package main

import "fmt"

// Person defines fields for an individual, such as name, age, and snn, etc.
type Person struct {
	name string
	age  int
	snn  string
	// TODO: 1 - add additional fields for a person's address
	addr string
}

func main() {
	// TODO: 2 - create a slice of Person containing 3 persons

	p := []Person{
		{name: "张三", age: 16, snn: "1-1-1", addr: "北京"},
		{name: "李四", age: 17, snn: "1-1-2", addr: "上海"},
		{name: "王五", age: 18, snn: "1-1-3", addr: "广州"},
	}

	// print each person's name and city/state ONLY
	for i, p := range p {
		fmt.Printf("%v. %v住在%v\n",
			i, p.name, p.addr)
	}

	p1 := Person{name: "小青", age: 19, snn: "1-1-3", addr: "深圳"}

	p = append(p, p1)

	for _, p := range p {
		fmt.Printf("%v,年龄%v,住在%v\n", p.name, p.age, p.addr)
	}
}
