package decoder

import "bufio"

// 基于分隔符的拆包器

type DelimiterDecoder struct {
	delimiter byte
	rd        *bufio.Reader
}

func NewDelimiterDecoder(delimiter byte) Decoder {
	return &DelimiterDecoder{
		delimiter: delimiter,
	}
}
func (d *DelimiterDecoder) ReadFrame() ([]byte, error) {
	payload, err := d.rd.ReadSlice(d.delimiter)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
func (d *DelimiterDecoder) BindReader(rd *bufio.Reader) {
	d.rd = rd
}
