package guochan

import (
	"crypto/aes"
	"crypto/cipher"
	"hash"
	"io"

	"golang.org/x/crypto/hkdf"
)

// AEADWriter is ...
type AEADWriter struct {
	A     cipher.AEAD
	W     io.WriteCloser
	nonce []byte
}

func (w AEADWriter) Write(src []byte) (n int, err error) {
	c := make([]byte, len(src)+w.A.Overhead())
	c = w.A.Seal(nil, w.nonce, src, nil)
	n, err = w.W.Write(c)
	if n != 0 && err != nil && err != io.EOF {
		return 0, err
	}

	if n == 0 {
		return 0, err
	}

	return len(src), nil
}

// Close is ...
func (w AEADWriter) Close() error {
	return w.W.Close()
}

// AEADReader is ...
type AEADReader struct {
	A     cipher.AEAD
	R     io.Reader
	nonce []byte
}

func (r AEADReader) Read(dst []byte) (n int, err error) {
	n, err = r.R.Read(dst)
	if n > 0 && err != nil && err != io.EOF {
		return 0, err
	}

	r.A.Open(dst[:0], r.nonce, dst[:n], nil)
	if err != nil {
		return 0, err
	}
	return n - r.A.Overhead(), nil
}

// NewGCMWithHKDF is ...
func NewGCMWithHKDF(hash func() hash.Hash, secret, salt, info []byte) (cipher.AEAD, error) {

	hkdfrd := hkdf.New(hash, secret, salt, info)

	key := make([]byte, aes.BlockSize)

	_, err := io.ReadFull(hkdfrd, key)

	if err != nil {
		return nil, err
	}

	blk, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(blk)

}
