package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

var (
	sarRe   = regexp.MustCompile(`Average:.*`)
	idleRe  = regexp.MustCompile(`\d+\.\d+`)
	totalRe = regexp.MustCompile(`MemTotal:\s+(\d+)`)
	freeRe  = regexp.MustCompile(`MemFree:\s+(\d+)`)
)

func main() {

	var idle float64

	flag.Float64Var(&idle, "c", 80, "cpu idle:10 - 90")
	flag.Parse()

	if idle < 10 || idle > 90 {
		log.Fatalln("cpu idle must idle >= 10 && idle <= 90")
	}

	ticker := time.NewTicker(2 * time.Second)

	t := new(tickers)

	m := new(mem)

	memTotal, err := getMemTotal()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for range ticker.C {
			m.getMemFree()
			fmt.Println("free: ", m.free)
			if m.free/memTotal > 0.5 {
				m.allocate()
			} else {
				m.freemem()
			}
		}
	}()

	for range ticker.C {

		t.cpuidle()
		fmt.Println("idle: ", t.idle)
		if t.idle > idle {
			t.push()
		} else {
			t.pop()
		}
		fmt.Println("tickers: ", runtime.NumGoroutine())
	}

}

func getMemTotal() (float64, error) {

	data, err := ioutil.ReadFile("/proc/meminfo")

	if err != nil {
		return 0, err
	}

	total := totalRe.FindSubmatch(data)[1]
	totalf, err := strconv.ParseFloat(string(total), 64)

	if err != nil {
		return 0, err
	}

	return totalf, nil

}

type tickers struct {
	ch   []chan struct{}
	idle float64
}

type mem struct {
	bufs []*[]byte
	free float64
}

func (m *mem) getMemFree() {
	data, err := ioutil.ReadFile("/proc/meminfo")

	if err != nil {
		m.free = 0
		log.Println(err)
		return
	}

	free := freeRe.FindSubmatch(data)[1]
	freef, err := strconv.ParseFloat(string(free), 64)

	if err != nil {
		m.free = 0
		log.Println(err)
		return
	}

	m.free = freef

}

func (m *mem) allocate() {
	if m.free == 0 {
		return
	}
	b := make([]byte, int(m.free)*300)
	for k, _ := range b {
		b[k] = '0'
	}

	m.bufs = append(m.bufs, &b)
}

func (m *mem) freemem() {
	l := len(m.bufs)

	if l == 0 {
		return
	}

	c := m.bufs[l-1]
	*c = nil
	c = nil
	m.bufs = m.bufs[:l-1]
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
	if numcpu := runtime.NumCPU(); numcpu > 16 {
		n = numcpu
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
