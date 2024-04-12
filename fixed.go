package decoder

import "bufio"

// 定长拆包器
type FixedDecoder struct {
	// 拆包器大小
	size   int
	reader *bufio.Reader
}

func NewFixedDecoder(size int, reader *bufio.Reader) Decoder {
	return &FixedDecoder{size: size, reader: reader}
}
func (d *FixedDecoder) ReadFrame() ([]byte, error) {
	var data = make([]byte, d.size)
	if _, err := d.reader.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
func (d *FixedDecoder) BindReader(rd *bufio.Reader) {
	d.reader = rd
}
