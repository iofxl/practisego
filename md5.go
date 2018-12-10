package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	fmt.Println("key is:", kdf("foo", 24))
}

// b = h.Sum(b) 会把hash结果加到b上去，返回最终的总的，prev就是上一次的哈希结果，所以这个算法就是把上一次的哈希结果跟password 一时反复的计算，到积累够要的长度
func kdf(password string, keyLen int) []byte {
	var b, prev []byte
	h := md5.New()
	for len(b) < keyLen {

		fmt.Println("b is:", b)
		fmt.Println("prev is:", prev)

		h.Write(prev)
		h.Write([]byte(password))

		b = h.Sum(b)

		prev = b[len(b)-h.Size():]
		h.Reset()
	}
	return b[:keyLen]
}
