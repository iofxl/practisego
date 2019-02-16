package guochan

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"log"
	"net"
	"time"
)

// G国产器 is G
type G国产器 interface {
	G国产(conn net.Conn) (net.Conn, error)
}

type listener struct {
	net.Listener
	G国产器
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return l.G国产(c)
}

// Listen is Listen
func Listen(g G国产器, network, address string) (net.Listener, error) {
	l, err := net.Listen(network, address)
	return &listener{l, g}, err
}

// DialTimeout is ...
func DialTimeout(g G国产器, network, address string, timeout time.Duration) (net.Conn, error) {
	c, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}

	return g.G国产(c)
}

// Method is Method
type Method uint

const (
	none Method = iota
	// Aes128Ctr is AES_128_CTR
	Aes128Ctr
	// Aes192Ctr is AES_192_CTR
	Aes192Ctr
	// Aes256Crt is Aes256Crt
	Aes256Crt
	maxMethod
)

var keySizes = []uint8{
	none:      0,
	Aes128Ctr: 16,
	Aes192Ctr: 24,
	Aes256Crt: 32,
}

type secret struct {
	m Method
	s []byte
}

var s secret

func kdf(size int, secret []byte) []byte {
	var prev, key []byte

	h := md5.New()

	for len(key) < size {

		h.Write(prev)
		h.Write(secret)
		// Sum appends the current hash to b  and returns the resulting slice.
		key = h.Sum(key)

		prev = key[len(key)-h.Size():]
		h.Reset()
	}

	return key[:size]
}

func (s *secret) newCipher() cipher.Block {

	key := kdf(int(keySizes[s.m]), s.s)

	blk, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	return blk
}

func (s *secret) G国产(c net.Conn) (net.Conn, error) {
	// so the work MUST be done in NewStreamConn
	return NewStreamConn(c, s.newCipher())
}

// New国产器 is New国产器
func New国产器(m Method, s []byte) G国产器 {
	return &secret{m, s}
}

/*
var methods = make([]func() 国产器, maxMethod)

func Register国产器(m Method, f func() 国产器) {

	if m >= maxMethod {

		panic("国产: Register国产器 of unknown guochan function")

	}

	methods[m] = f

}

func (m Method) Size() int {
	if h > 0 && h < maxMethod {
		return int(keySizes[m])
	}
	panic("国产: Size of unknown method function")
}

// 到底是否需要每一种方法去实现这个接口呢？还是做成每一种方法有能力New出一个这
// 种接口更好？我感觉是有法New出更好？那么为什么感觉更好呢？

// 但是这里我感觉得以这个m去实现这个接口非常合适，这是为什么呢？而还是使hash方
// 法的那套，先注册再使New去执行

// conn类型只包含到Block Interface是非常合适的，因为conn最多只接受到Block了过，
// 这是它需要的东西，神使conn包含secret的曰事，每一次新连接都要计算一遍blk

// 这里问题里罢，哪样使m 跟 blk 产生联系呢？还有使得一个secret参数

// m 跟 blk, 只有blk 跟目标net.Conn有联系，要从这里入手, 所以是必顉把blk包装一
// 下再去实现转换接口？

// 瞅以方法是无资格来实现这个接口的，这个接口只能有密码来实现，实现接口的东西肯
// 定是需要作为参数的东西
func (m Method) 国产(c net.Conn) net.Conn {
	return &conn{Conn: c, blk: blk}
}
*/
