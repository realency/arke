package bits

import (
	"io"
)

type ArrayReader struct {
	bytes  []byte
	index  int
	offset int
	mask   byte
	right  byte
}

func NewArrayReader(bytes []byte, bitOffset int) *ArrayReader {
	if bitOffset < 0 || bitOffset > len(bytes)*8 {
		panic("arg out of range")
	}

	offset := bitOffset % 8
	index := bitOffset / 8

	return &ArrayReader{
		bytes:  bytes,
		index:  index,
		offset: offset,
		mask:   0xFF >> offset,
		right:  bytes[index] << offset,
	}
}

func (r *ArrayReader) Read(p []byte) (int, error) {
	if r.index >= len(r.bytes) {
		return 0x00, io.EOF
	}

	// optimise for the common case of a zero offset
	if r.offset == 0 {
		count := copy(p, r.bytes[r.index:])
		r.index += count
		return count, nil
	}

	var left byte = 0x00
	i := 0

	for r.index < len(r.bytes)-1 && i < len(p) {
		r.index++
		b := r.bytes[r.index]
		left = b >> (8 - r.offset)
		p[i] = r.right | left
		r.right = b << r.offset
		i++
	}

	if r.index == len(r.bytes)-1 && i < len(p) {
		p[i] = r.right
		r.index++
		i++
	}

	return i, nil
}
