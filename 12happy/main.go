package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	mathrand "math/rand"
	"time"
)

func main() {

	done := make(chan struct{})

	for i := 0; i < 100; i++ {
		go work(done)
	}

	<-done

}

func work(done chan struct{}) {

	go func() {

		key := make([]byte, aes.BlockSize)
		iv := make([]byte, aes.BlockSize)
		rand.Read(key)
		rand.Read(iv)
		blk, _ := aes.NewCipher(key)
		stream := cipher.NewCTR(blk, iv)

		src := make([]byte, 4096*4096)
		rand.Read(src)
		dst := make([]byte, len(src))

		for {
			stream.XORKeyStream(dst, src)
			stream.XORKeyStream(src, dst)
			time.Sleep(time.Duration(mathrand.Intn(8)) * time.Second)
			fmt.Println("next")
		}
		close(done)

	}()

}
