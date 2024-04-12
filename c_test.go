package decoder

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"testing"
	"time"
)

type MockConn struct {
	net.Conn
	data string
}

func (c *MockConn) Read(b []byte) (n int, err error) {
	result, err := hex.DecodeString(c.data)
	if err != nil {
		return 0, err
	}
	copy(b, result)
	// mock dely
	time.Sleep(10 * time.Millisecond)
	return len(result), nil
}

func TestLengthFileConninitialBytesToStrip2(t *testing.T) {
	data := "01032c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000007859"
	content := ""
	for i := 0; i < 10; i++ {
		content += data
	}

	c := NewConn(&MockConn{data: content}, NewLengthFieldBasedDecoder(2, 65535, 1, 2, 0, binary.BigEndian), newEncoder())
	for i := 0; i < 10; i++ {
		f, err := c.ReadFrame()
		if err != nil {
			t.Fatalf(err.Error())
		}

		if hex.EncodeToString(f) != data {
			t.Fatalf("got [%s] but want [%s] ", hex.EncodeToString(f), data)
		}
	}

}
