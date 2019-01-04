package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"runtime/debug"
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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	mticker := time.NewTicker(5 * time.Second)
	defer mticker.Stop()

	t := new(tickers)

	go func() {
		var m mem
		for range mticker.C {
			t, f, err := getMemStats()
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Println("free: ", f)
			if f/t > 0.5 {
				m.add(int(f/5) * 1000)
			} else {
				m.free()
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

type tickers struct {
	ch   []chan struct{}
	idle float64
}

func getMemStats() (float64, float64, error) {
	data, err := ioutil.ReadFile("/proc/meminfo")

	if err != nil {
		return 0, 0, err
	}

	free := freeRe.FindSubmatch(data)[1]
	freef, err := strconv.ParseFloat(string(free), 64)

	if err != nil {
		return 0, 0, err
	}

	total := totalRe.FindSubmatch(data)[1]
	totalf, err := strconv.ParseFloat(string(total), 64)

	if err != nil {
		return 0, 0, err
	}

	return totalf, freef, nil

}

type mem [][]byte

func allocate(n int) []byte {
	b := make([]byte, n)
	for k, _ := range b {
		b[k] = '0'
	}

	return b

}

func (m *mem) add(n int) {

	*m = append(*m, allocate(n))

}

func (m *mem) free() {

	l := len(*m)

	if l == 0 {
		return
	}

	(*m)[l-1] = nil
	(*m) = (*m)[:l-1]
	debug.FreeOSMemory()
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
