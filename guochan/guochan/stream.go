package guochan

import (
	"crypto/cipher"
	"net"
)

// goss2的做法是通过AESCTR函数生成blk, 把blk使ctrStream结构包起,这个ctrStream实
// 现了这个包里的Cipher接口,这个Cipher接口的作用,也就是生成加密的工作模式, 瞅以
// 感觉相当乱

// 需要的是一个使net.Conn发生变化的接口，conn实现变化过的net.Conn接口，以
// net.Conn为参数，返回这个东西的函数事实就是实现这个接口

// 实际的东西的肯定要有能够区分的属性的。 而要做到上层有法使统一的变量或者函数
// 表示出来，参考sha256包 肯定有一个函数在这个变量里做文章, 从哪个位置开始分隔
// ？ 就是要从开始不同的位置 所有东西有法相同的东西 要尽量保持相同
// 各个功能模式关心的域不同，所以中间肯定要有一个函数做这种连接工作，作为功能之
// 间的接口

// goss2的做法是设计一个生成模式的接口，这个接口的名字就是Cipher, 直接把这个接
// 口包含在这个结构里面的(这个要包进去干嘛？)，而实现这个Cipher接口的东西包含密
// 码，包含一个生成模式 的函数(统一名字的作用，比如统一喊New)，包含这个函数就是
// 多余的，再去实现一个有法生成模式的接口，这种做法么好呢？ 么需要呢？

// goss2的做法，顶层设计一个加密接口，接口的作用是转换net.Conn，使一个只包含下
// 面那个接口的东西实现 不加密，stream, aead 三种加密方法，会实现自己的一种加密
// 接口，这个接口的作用是生成模式，使一个只包含块接口的东西实现

// 使一个包含net.Conn接口，和 上面这个加密接口的东西 去实现转换后的net.Conn接口
// 使一个函数返回上面要使的块接口，这个函数以key为参数, 使func PickCipher(
// name, password ) ( Cipher, error ) 返回这个转换接口，这个函数会调用各个方法
// 的New函数

// 本身包里每一种加密方法都有相应的一个结构实现这个接口，而这个结构其实就是包含
// 各种方法自己的一个同名接口，这个同名接口的作用是生成模式，这种弄一下有什么意
// 义呢？ 为什么不各种方法直接实现这个接口？ 我瞅以是多余的

// 而各个包里的同名接口其实就是一个返回加密模式的接口，然后又使包含一个块接口的
// 结构去实现，本身块就是直接有法生成模式接口的，那这种弄一下有什么意义呢？

// 首先这个接口的设计是很需要的，但是这个接口确实不是应该conn来实现的，conn要做
// 的是实现修改过的net.Conn接口，这个接口的作用就是一种描述，由加密方法本身去实
// 现就会得罢 我想，一开始的依据统都“我的问题的是什么，我要哪样去解决〞，但是这
// 样做的曰事，设计在哪里

// secret是需要外界传进来的，不需要留含进去呢？我的理解是不需要的，只需要包含真
// 正不同的变量就会得罢

// 我就是理解不进去，为什么goss2通过那种套一下又套一下，最终就做成了呢，我这种
// 做法的问题在哪里？

// 他的conn包含的是Cipher接口，事实就是包含的Block接口，无区别的这里

// StreamConn is ...
type StreamConn struct {
	net.Conn
	r *cipher.StreamReader
	w *cipher.StreamWriter
}

// 犯了个傻逼错误, 没有用指针去实现,结果总是已经initWriter()的conn, 要写的时候
// 还是报c.w是nil
// 我认为这种方法是不对的，应该在使NewConn的方式实现
/*
func (c *conn) initReader() error {
	if c.r == nil {
		iv := make([]byte, c.blk.BlockSize())
		if _, err := io.ReadFull(c.Conn, iv); err != nil {
			return err
		}
		ctr := cipher.NewCTR(c.blk, iv)
		c.r = &cipher.StreamReader{S: ctr, R: c.Conn}
	}
	return nil
}
*/

// NewStreamConn is ...
// the way use initReader and initWriter in Reader and Writer is funny.
func NewStreamConn(c net.Conn, blk cipher.Block) (net.Conn, error) {
	iv := make([]byte, blk.BlockSize())
	ctr := cipher.NewCTR(blk, iv)

	_, err := c.Write(iv)
	if err != nil {
		return nil, err
	}

	stream := &StreamConn{
		Conn: c,
		w:    &cipher.StreamWriter{S: ctr, W: c},
	}

	n, err := c.Read(iv)

	if err != nil {
		return nil, err
	}

	ctr = cipher.NewCTR(blk, iv[:n])

	stream.r = &cipher.StreamReader{S: ctr, R: c}

	return stream, nil

}

func (c *StreamConn) Read(p []byte) (int, error) {
	return c.r.Read(p)
}

// 比如，这里初始化的时候，我的问题就是我要生成模式，我需要有办法得到块变量，我
// 需要知道iv的长度，对于conn这个东西来曰，除了自带这两个信息，它还有别的什么方
// 法？ 或者自带另外一个有办法生成这两个信息的东西，那么哪样选择呢？
// iv的长度跟blk的长度一样
/*
func (c *conn) initWriter() error {
	if c.w == nil {
		iv := make([]byte, c.blk.BlockSize())
		_, err := io.ReadFull(rand.Reader, iv)
		if err != nil {
			return err
		}
		if _, err := c.Conn.Write(iv); err != nil {
			return err
		}
		ctr := cipher.NewCTR(c.blk, iv)
		c.w = &cipher.StreamWriter{S: ctr, W: c.Conn}

	}
	return nil
}
*/

// Write is ...
func (c *StreamConn) Write(p []byte) (int, error) {
	return c.w.Write(p)
}
