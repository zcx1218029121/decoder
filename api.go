package decoder

import (
	"bufio"
	"net"
)

type Decoder interface {
	// 拆包器 拆包返回完整数据包
	ReadFrame() ([]byte, error)
	BindReader(rd *bufio.Reader)
}
type Encoder interface {
	WriteFrame([]byte) error
	BindWriter(wr *bufio.Writer)
}
type Conn interface {
	net.Conn
	Decoder
	Encoder
	Close() error
	Flush() error
	SetDecoder(Decoder)
	SetEncoder(Encoder)
}
