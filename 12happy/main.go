package main

import "time"

func main() {

	// tick := time.NewTicker(time.Nanosecond) // 300%
	// tick := time.NewTicker(time.Microsecond) // 190%
	tick := time.NewTicker(time.Millisecond) // 9%
	// tick := time.NewTicker(time.Second) // 0%

	for range tick.C {
	}

}
