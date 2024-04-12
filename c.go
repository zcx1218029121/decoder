package decoder

import (
	"bufio"
	"net"
)

type ConnImpl struct {
	net.Conn
	Decoder
	Encoder
	rd *bufio.Reader
	wr *bufio.Writer
}

func (c *ConnImpl) SetDecoder(d Decoder) {
	d.BindReader(c.rd)
	c.Decoder = d
}
func (c *ConnImpl) SetEncoder(e Encoder) {
	e.BindWriter(c.wr)
	c.Encoder = e
}
func NewConn(conn net.Conn, decoder Decoder, encoder Encoder) Conn {
	c := &ConnImpl{
		Conn: conn,
		rd:   bufio.NewReaderSize(conn, 4096),
		wr:   bufio.NewWriterSize(conn, 1024),
	}
	c.SetDecoder(decoder)
	c.SetEncoder(encoder)
	return c
}
func (c *ConnImpl) Flush() error {
	return c.wr.Flush()
}
func NewConnWithRw(conn net.Conn, decoder Decoder, encoder Encoder, rd *bufio.Reader, wr *bufio.Writer) Conn {
	c := &ConnImpl{
		Conn: conn,
		rd:   rd,
		wr:   wr,
	}
	c.SetDecoder(decoder)
	c.SetEncoder(encoder)
	return c
}
