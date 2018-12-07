package main

// 你不能賦值給一个nil指针

import (
	"fmt"
	"strconv"
	"sync"
)

type node struct {
	data int
	next *node
}

type LinkedList struct {
	head *node
	size int
}

type splitResult struct {
	prev, head, tail *node
}

// 切分要注意的:
// 1. 要先得到下一个的头, 再传出去
func splitList(ll *LinkedList, subLen int) <-chan splitResult {

	if subLen <= 0 || subLen > ll.size {
		subLen = ll.size
	}

	ch := make(chan splitResult)

	go func() {

		index := ll.head
		head := index
		prev := new(node)
		for index != nil {

			result := splitResult{new(node), new(node), new(node)}

			for size := 1; size < subLen && index.next != nil; size++ {
				index = index.next
			}

			if result.head == ll.head {
				result.prev = nil
			} else {
				result.prev = prev
			}

			result.head = head
			result.tail = index

			prev = index
			index = index.next
			head = index
			ch <- result
		}

		close(ch)

	}()

	return ch
}

func worker(wg *sync.WaitGroup, resultch <-chan splitResult, head chan *node) error {

	defer wg.Done()

	for result := range resultch {

		if result.head == result.tail {
			if result.tail.next != nil {
				result.head.next = result.prev
				return nil
			} else {
				newnode := new(node)
				newnode.data = result.tail.data
				result.tail = newnode
				result.tail.next = result.prev
				head <- result.tail
				return nil
			}
		}

		prev, move, move1 := result.head, result.head.next, result.head.next
		for move != result.tail {
			move1 = move
			move = move.next
			move1.next = prev
			prev = move1

		}

		// 不能給nil的prev賦值.其实也就是头会发生这样的情况
		result.head.next = result.prev

		// 总之无论如何不能給Nil指针賦值
		if result.tail.next == nil {
			newnode := new(node)
			newnode.data = result.tail.data
			result.tail = newnode
			result.tail.next = prev
			head <- result.tail
		} else {
			result.tail.next = prev
		}

	}

	return nil

}

func (ll *LinkedList) Reverse(subLen int, count int) error {

	var wg sync.WaitGroup

	headch := make(chan *node)

	splitCh := splitList(ll, subLen)

	// you can run any GR
	for i := 0; i < count; i++ {
		wg.Add(1)
		go worker(&wg, splitCh, headch)
	}

	ll.head = <-headch
	wg.Wait()

	return nil
}

// Insert 是要注意的,插入的第一个节点的next指针是不指定的.
func (ll *LinkedList) Insert(e ...int) int {

	for _, data := range e {

		newNode := new(node)
		newNode.data = data

		if ll.size == 0 {
			ll.head = newNode
		} else {
			newNode.next = ll.head
			ll.head = newNode
		}
		ll.size++
	}

	return ll.size
}

func (ll *LinkedList) String() string {

	var b []byte

	b = append(b, []byte("LinkedList:\n")...)

	index := ll.head
	for i := 0; i < ll.size; i++ {
		b = strconv.AppendInt(b, int64(index.data), 10)
		b = append(b, ' ')
		index = index.next
	}
	b = append(b, '\n')

	return string(b)
}

func main() {

	ll := new(LinkedList)

	//	data := []int{19, 30, 7, 23, 24, 0, 12, 28, 3, 11, 18, 1, 31, 14, 21, 2, 9, 16, 4, 26, 10, 25}

	N := 50
	data := make([]int, N)
	for i := 0; i < N; i++ {
		data[i] = i
	}

	ll.Insert(data...)
	fmt.Println(ll)
	ll.Reverse(63, 1)
	fmt.Println(ll)

}
