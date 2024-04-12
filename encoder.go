package decoder

import "bufio"

// 简单实现的的encoder 数据直接写入 不暴露给外界调用
type simpleEncoder struct {
	wr *bufio.Writer
}

func (e *simpleEncoder) WriteFrame(data []byte) error {
	_, err := e.wr.Write(data)
	return err
}
func newEncoder() Encoder {
	return &simpleEncoder{}
}
func (d *simpleEncoder) BindWriter(wr *bufio.Writer) {
	d.wr = wr
}

