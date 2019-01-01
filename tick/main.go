package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

var (
	sarRe  = regexp.MustCompile(`Average:.*`)
	idleRe = regexp.MustCompile(`\d+\.\d+`)
)

func main() {

	ticker := time.NewTicker(2 * time.Second)

	t := new(tickers)

	for range ticker.C {

		t.cpuidle()
		fmt.Println("idle: ", t.idle)
		if t.idle > 80 {
			t.push()
		} else {
			t.pop()
		}
		fmt.Println("tickers: ", runtime.NumGoroutine())
	}

}

type tickers struct {
	ch   []chan struct{}
	idle float64
}

func (t *tickers) cpuidle() {

	out, err := exec.Command("sar", "-u", "2", "2").CombinedOutput()
	if err != nil {
		t.idle = 0
		log.Println(out, err)
		return
	}

	m := idleRe.FindAll(sarRe.Find(out), -1)
	f, err := strconv.ParseFloat(string(m[5]), 64)

	if err != nil {
		t.idle = 0
		log.Println(err)
		return
	}

	t.idle = f
}

func (t *tickers) push() {
	stop := make(chan struct{})
	n := 16
	if runtime.NumCPU() > 16 {
		n = runtime.NumCPU()
	}
	for i := 0; i < n; i++ {
		go tick(stop)
	}
	t.ch = append(t.ch, stop)
}

func (t *tickers) pop() {
	l := len(t.ch)
	if l < 1 {
		return
	}
	close(t.ch[l-1])
	t.ch = t.ch[:l-1]
}

func tick(stop chan struct{}) {
	go func() {
		t := time.NewTicker(time.Millisecond)
		for range t.C {
			select {
			case <-stop:
				return
			default:
			}
		}
	}()
}
