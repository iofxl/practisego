package main

import (
	"fmt"
	"runtime"
	"time"

	"bitbucket.org/bertimus9/systemstat"
)

type ticker struct {
	ch     []chan struct{}
	idle   chan float64
	want   uint
	numcpu int
}

func (t *ticker) updateIdle() {

	go func() {

		tick := time.NewTicker(2 * time.Second)
		defer tick.Stop()

		lastsmp := systemstat.GetCPUSample()

		for range tick.C {

			cur := systemstat.GetCPUSample()

			t.idle <- float64(cur.Idle-lastsmp.Idle) / float64(cur.Total-lastsmp.Total) * 100
			lastsmp = cur
		}
	}()

}

func (t *ticker) push() {

	stop := make(chan struct{})

	for i := 0; i < t.numcpu; i++ {
		go tick(stop)
	}

	t.ch = append(t.ch, stop)
}

func (t *ticker) pop() {
	l := len(t.ch)
	if l < 1 {
		return
	}

	close(t.ch[l-1])
	t.ch = t.ch[:l-1]
}

func (t *ticker) clear() {
	l := len(t.ch)
	if l < 1 {
		return
	}

	for k, _ := range t.ch {
		close(t.ch[k])
	}

	t.ch = nil
}

func tick(stop chan struct{}) {
	go func() {
		t := time.NewTicker(time.Millisecond)
		defer t.Stop()
		for range t.C {
			select {
			case <-stop:
				return
			default:
			}
		}
	}()
}

func (t *ticker) run() {

	t.updateIdle()

	for v := range t.idle {
		fmt.Println("cpu used: ", 100-v)
		fmt.Println("tickers: ", runtime.NumGoroutine())

		if 100-v < float64(t.want) {
			t.push()
			continue
		}

		if 100-v >= float64(t.want+20) {
			t.clear()
			time.Sleep(5 * time.Minute)
		}

	}

}
