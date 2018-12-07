package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/hkdf"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s passwd \n", os.Args[0])
		os.Exit(1)
	}

	passwd := os.Args[1]

	hash := sha1.New

	secret := []byte(passwd)

	salt := make([]byte, hash().Size())

	_, err := rand.Read(salt)

	if err != nil {
		log.Fatal(err)
	}

	info := []byte("ss-subkey")

	hf := hkdf.New(hash, secret, salt, info)

	key := make([]byte, 16)

	_, err = hf.Read(key)

	if err != nil {
		log.Fatal(err)
	}

	for {
		var text []byte

		fmt.Scanf("%v\n", &text)

		if n := len(text) % aes.BlockSize; n != 0 {

			n = aes.BlockSize - n

			log.Printf("n is: %v\n", n)

			t := make([]byte, len(text)+n)
			copy(t, text)
			text = t

			log.Printf("len text now is: %v\n", len(text))

			// "plaintext is not a multiple of the block size"
		}

		block, err := aes.NewCipher([]byte(key))

		if err != nil {
			log.Println("Generate blk error", err)
			continue
		}

		// can't just var ciphertext []byte here, must make it
		ciphertext := make([]byte, aes.BlockSize+len(text))

		iv := ciphertext[:aes.BlockSize]

		_, err = rand.Read(iv)

		if err != nil {
			log.Fatal(err)
		}

		enc := cipher.NewCBCEncrypter(block, iv)
		dec := cipher.NewCBCDecrypter(block, iv)

		enc.CryptBlocks(ciphertext[aes.BlockSize:], text)

		fmt.Println("ciphertext:", hex.EncodeToString(ciphertext))

		newtext := make([]byte, len(text))
		dec.CryptBlocks(newtext, ciphertext[16:])

		fmt.Println("plaintext:", string(newtext))
	}

}
