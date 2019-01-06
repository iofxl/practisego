package main

import (
	"fmt"
	"runtime/debug"
	"time"

	"bitbucket.org/bertimus9/systemstat"
)

type mem struct {
	bufs     [][]byte
	used     chan float64
	memFree  uint64
	want     uint
	pagesize int
}

func (m *mem) allocate(n int) {

	b := make([]byte, n)

	for i := 0; i < n; i += m.pagesize {
		b[i] = 1
	}

	m.bufs = append(m.bufs, b)

}

func (m *mem) free() {

	l := len(m.bufs)

	if l == 0 {
		return
	}

	m.bufs[l-1] = nil
	m.bufs = m.bufs[:l-1]
	debug.FreeOSMemory()
}

func (m *mem) freeAll() {

	l := len(m.bufs)

	if l == 0 {
		return
	}

	for i := 0; i < l; i++ {
		m.bufs[i] = nil
	}
	m.bufs = m.bufs[:0]
	debug.FreeOSMemory()
}

func (m *mem) updateMemUsed() {

	go func() {

		t := time.NewTicker(2 * time.Second)
		defer t.Stop()

		for range t.C {

			cur := systemstat.GetMemSample()

			m.memFree = cur.MemFree
			m.used <- float64(cur.MemUsed) / float64(cur.MemTotal) * 100
		}
	}()
}

func (m *mem) run() {

	m.updateMemUsed()

	for v := range m.used {
		fmt.Println("mem used: ", v)

		if v < float64(m.want) {
			m.allocate(int(m.memFree) * 100)
			continue
		}

		if v >= float64(m.want)+10 {
			m.freeAll()
			time.Sleep(5 * time.Minute)
		}

	}

}
