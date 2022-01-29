package bits

import "io"

type sliceReader struct {
	bytes  []byte
	index  int
	offset int
	ready  byte
}

func newReader(bytes []byte, bitOffset int) *sliceReader {
	if bitOffset < 0 || bitOffset > len(bytes)*8 {
		panic("arg out of range")
	}

	offset := bitOffset % 8
	index := bitOffset / 8

	return &sliceReader{
		bytes:  bytes,
		index:  index,
		offset: offset,
		ready:  bytes[index] >> offset,
	}
}

func (r *sliceReader) Read(p []byte) (int, error) {
	if r.index >= len(r.bytes) {
		return 0, io.EOF
	}

	// optimise for the common case of a zero offset
	if r.offset == 0 {
		count := copy(p, r.bytes[r.index:])
		r.index += count
		return count, nil
	}

	i := 0
	for r.index < len(r.bytes)-1 && i < len(p) {
		r.index++
		b := r.bytes[r.index]
		p[i] = r.ready | (b << (8 - r.offset))
		r.ready = b >> r.offset
		i++
	}

	if r.index == len(r.bytes)-1 && i < len(p) {
		p[i] = r.ready
		r.index++
		i++
	}

	return i, nil
}
