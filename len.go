package decoder

import (
	"bufio"
	"encoding/binary"
	"errors"
)

// 基于长度域的拆包器
var ErrOverMax = errors.New("over max len")

func (d *LengthFieldBasedDecoder) BindReader(rd *bufio.Reader) {
	d.rd = rd
}

type LengthFieldBasedDecoder struct {
	rd *bufio.Reader
	//指定长度字段的偏移量，
	//也就是长度字段在整个数据包中位于哪个位置，以字节为单位。
	//例如，如果长度字段是包头的第3个字节，则偏移量应该是2
	lengthFieldOffset uint64
	//指定最大的帧长度。如果超出此长度，则抛出一个异常并丢弃该帧。。
	maxFrameLength uint64
	//指定从长度字段指示的长度中需要添加或减去的额外字节数。例如，如果数据包总长度是包含长度字段本身的，则长度调整值应该为负的长度字段长度
	lengthAdjustment int64
	//指定跳过多少字节才能开始解码有效数据。例如，可以设置为长度字段自己的长度，因为在得到长度后，它不再需要被解码。
	initialBytesToStrip uint64
	//长度域长度 当位数为 3,5,6,7的时候会向上对齐
	lengthFieldLength uint64
	// 长度域大小端
	byteOrder binary.ByteOrder
}

func (d *LengthFieldBasedDecoder) ReadFrame() ([]byte, error) {
	var frame []byte
	// 跳过指定长度
	var offset = make([]byte, d.lengthFieldOffset)
	_, err := d.rd.Read(offset)
	if err != nil {
		return nil, err
	}
	var lenFiled = make([]byte, d.lengthFieldLength)
	_, err = d.rd.Read(lenFiled)
	if err != nil {
		return nil, err
	}
	// 对长度域进行对齐
	var payload []byte
	switch d.lengthFieldLength {
	case 1:
		{
			var actLen = int64(uint8(lenFiled[0])) + d.lengthAdjustment
			if actLen < 0 {
				return nil, ErrOverMax
			}
			if actLen > int64(d.maxFrameLength) {
				return nil, ErrOverMax
			}
			payload = make([]byte, actLen)
			_, err := d.rd.Read(payload)
			if err != nil {
				return nil, err
			}

		}
	case 2:
		{
			var actLen = int64(d.byteOrder.Uint16(lenFiled)) + d.lengthAdjustment
			if actLen < 0 {
				return nil, ErrOverMax
			}
			if actLen > int64(d.maxFrameLength) {
				return nil, ErrOverMax
			}
			payload = make([]byte, actLen)
			_, err := d.rd.Read(payload)
			if err != nil {
				return nil, err
			}
		}
	case 3, 4:
		{
			var expLenFiled = alignment(d.byteOrder, lenFiled, 4)
			var actLen = int64(d.byteOrder.Uint32(expLenFiled)) + d.lengthAdjustment
			if actLen < 0 {
				return nil, ErrOverMax
			}
			if actLen > int64(d.maxFrameLength) {
				return nil, ErrOverMax
			}
			payload = make([]byte, actLen)
			_, err := d.rd.Read(payload)
			if err != nil {
				return nil, err
			}

		}
	case 5, 6, 7, 8:
		{

			var expLenFiled = alignment(d.byteOrder, lenFiled, 8)
			var actLen = int64(d.byteOrder.Uint64(expLenFiled)) + d.lengthAdjustment
			if actLen < 0 {
				return nil, ErrOverMax
			}
			if actLen > int64(d.maxFrameLength) {
				return nil, ErrOverMax
			}
			payload = make([]byte, actLen)
			_, err := d.rd.Read(payload)
			if err != nil {
				return nil, err
			}
		}
	}
	frame = append(frame, offset...)
	frame = append(frame, lenFiled...)
	frame = append(frame, payload...)
	return frame[d.initialBytesToStrip:], nil
}

// 补0对齐数据长度
func alignment(order binary.ByteOrder, lenFiled []byte, exp uint64) []byte {

	if len(lenFiled) == int(exp) {
		return lenFiled
	}

	if len(lenFiled) > int(exp) {
		return lenFiled[:exp]
	}
	//大端前面补0
	var dst = make([]byte, exp)
	// 小端后面补0
	if order == binary.LittleEndian {
		copy(dst, lenFiled)
		return dst
	}
	// 大端前面补0
	copy(dst[len(dst)-len(lenFiled):], lenFiled)
	return dst
}
func NewLengthFieldBasedDecoder(
	lengthFieldOffset uint64,
	maxFrameLength uint64,
	lengthFieldLength uint64,
	lengthAdjustment int64,
	initialBytesToStrip uint64,
	byteOrder binary.ByteOrder) Decoder {
	return &LengthFieldBasedDecoder{
		lengthFieldOffset:   lengthFieldOffset,
		maxFrameLength:      maxFrameLength,
		lengthAdjustment:    lengthAdjustment,
		initialBytesToStrip: initialBytesToStrip,
		lengthFieldLength:   lengthFieldLength,
		byteOrder:           byteOrder,
	}
}
