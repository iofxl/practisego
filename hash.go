package main

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding"
	"fmt"
)

func main() {

	data := []byte("Hello World!")

	// func direct
	// these func return is arrary, Hash.Sum() return is slice.
	sum1 := sha256.Sum256(data)
	sum2 := sha256.Sum224(data)
	sum3 := sha512.Sum512(data)
	sum4 := sha512.Sum384(data)
	sum5 := sha512.Sum512_256(data)
	sum6 := sha512.Sum512_224(data)
	sum7 := sha1.Sum(data)
	sum8 := md5.Sum(data)

	// or New a digest first
	/*
		t: digest
		    d.I: Hash, BM, BU
	*/
	d1 := sha256.New()
	d2 := sha256.New224()
	d3 := sha512.New()
	d4 := sha512.New384()
	d5 := sha512.New512_256()
	d6 := sha512.New512_224()
	d7 := sha1.New()
	d8 := md5.New()

	d1.Write(data)
	d2.Write(data)
	d3.Write(data)
	d4.Write(data)
	d5.Write(data)
	d6.Write(data)
	d7.Write(data)
	d8.Write(data)

	hash1 := d1.Sum(nil)
	hash2 := d2.Sum(nil)
	hash3 := d3.Sum(nil)
	hash4 := d4.Sum(nil)
	hash5 := d5.Sum(nil)
	hash6 := d6.Sum(nil)
	hash7 := d7.Sum(nil)
	hash8 := d8.Sum(nil)

	print("use func:", "Size:32", sum1, "28", sum2, "64", sum3, "48", sum4, "32", sum5, "28", sum6, "20", sum7, "16", sum8)
	print("\nuse itf:", hash1, hash2, hash3, hash4, hash5, hash6, hash7, hash8)

	hashItf1 := crypto.SHA256.New()
	hashItf1.Write(data)
	hash9 := hashItf1.Sum(nil)

	print("\nuse type crypto.Hash:", hash9)

	state, _ := d1.(encoding.BinaryMarshaler).MarshalBinary()

	another := sha256.New()
	_ = another.(encoding.BinaryUnmarshaler).UnmarshalBinary(state)

	d1.Write([]byte("input2"))
	another.Write([]byte("input2"))

	hash1 = d1.Sum(nil)
	hashanother := another.Sum(nil)

	print("\nBinaryMarshaler ex.:", hash1, hashanother)

}

func print(e ...interface{}) error {
	for _, v := range e {

		if _, ok := v.(string); ok {
			_, err := fmt.Printf("%s\n", v)
			if err != nil {
				return err
			}

		} else {

			_, err := fmt.Printf("%x\n", v)
			if err != nil {
				return err
			}

		}

	}

	return nil
}
