package main

import (
	"testing"
)

func TestCalc(t *testing.T) {

	foo := map[int]int{
		0:     2,
		-1:    1,
		1:     3,
		99999: 100001,
	}

	for k, v := range foo {

		if output := Calc(k); output != v {
			t.Error("Test Failed: {} inputted, {} expected, recieved: {}", k, v, output)
		}

	}
}
