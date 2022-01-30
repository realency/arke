package bits

import (
	"io"
)

type writer struct {
	buff      []uint32
	index     int
	offset    int
	available int
	mask      uint32
	dest      uint32
}

func (w *writer) WriteByte(b byte) error {
	if w.available <= 0 {
		return io.EOF
	}

	if w.available < 8 {
		trunc := 8 - w.available
		b = (b >> trunc) << trunc
	}

	source := (uint32(b) << 24)
	w.dest = (w.dest & w.mask) | (source >> w.offset)

	w.buff[w.index] = w.dest

	w.offset += 8
	w.mask = (w.mask >> 8) | 0xFF000000

	if w.offset == 32 {
		w.available -= 8
		if w.available <= 0 {
			return nil
		}
		w.offset = 0
		w.index++
		w.mask = 0x00FFFFFF
		w.dest = w.buff[w.index]
	}

	if w.offset > 32 {
		w.offset -= 32
		w.available -= 8 - w.offset
		if w.available <= 0 {
			return nil
		}
		w.index++
		w.dest = w.buff[w.index]
		w.mask = (0xFFFFFFFF << w.offset) >> w.offset
		w.dest = (w.dest & w.mask) | (source << (8 - w.offset))
		w.mask = 0x00FFFFFF >> w.offset
		w.buff[w.index] = w.dest
	}

	return nil
}

func (m *Matrix) Writer(row, col, count int) io.ByteWriter {
	if row < 0 || row >= m.height || col < 0 || col >= m.width {
		panic("Arg out of bounds")
	}
	index := (row * m.intsPerRow) + (col / 32)
	offset := col % 32
	if m.width-col < count {
		count = m.width - col
	}
	return &writer{
		buff:      m.bits,
		index:     index,
		offset:    offset,
		available: count,
		dest:      m.bits[index],
		mask:      (0xFF000000 >> offset) ^ 0xFFFFFFFF,
	}
}
