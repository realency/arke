package bits

import "io"

type Direction byte

const (
	Up    Direction = 0x00
	Down  Direction = 0x01
	Left  Direction = 0x02
	Right Direction = 0x03
)

type reader struct {
	buff      []uint32  // The buffer - bits in bitfields as UInt32
	index     int       // index into the buffer of the current read position
	offset    int       // offset into the indexed uint of the current read position
	available int       // number of BITS available still to read
	direction Direction // Read direction - up, down, left, or right
	ready     uint32    // Pre-read bitfield for optimisation purposes - used in left- and right-readers, but not up- or down-.
	bit       uint32    // Mask to select the correct bit from the bitfield - used in up- and down-readers, but not left- or right-.
}

func (r *reader) readUp() byte {
	result := byte(0x0)

	var mask byte
	if r.available >= 8 {
		mask = 0x01
	} else {
		mask = 0x01 << (8 - r.available)
	}

	for i := 0; i < 8 && r.available > 0; i++ {
		result <<= 1
		if (r.buff[r.index] & r.bit) != 0 {
			result |= mask
		}
		r.index--
		r.available--
	}

	return result
}

func (r *reader) readDown() byte {
	result := byte(0x0)

	var mask byte
	if r.available >= 8 {
		mask = 0x01
	} else {
		mask = 0x01 << (8 - r.available)
	}

	for i := 0; i < 8 && r.available > 0; i++ {
		result <<= 1
		if (r.buff[r.index] & r.bit) != 0 {
			result |= mask
		}
		r.index++
		r.available--
	}

	return result
}

func (r *reader) readLeft() byte {
	result := byte(0x00)

	var mask byte
	if r.available >= 8 {
		mask = 0x01
	} else {
		mask = 0x01 << (8 - r.available)
	}

	for i := 0; i < 8 && r.available > 0; i++ {
		result <<= 1
		if (r.ready & 0x00000001) != 0 {
			result |= mask
		}
		r.ready >>= 1

		r.offset--
		if r.offset == -1 {
			r.offset = 31
			r.index--
			r.ready = r.buff[r.index]
		}
		r.available--
	}

	return result
}

func (r *reader) readRight() byte {
	result := byte(r.ready >> 24)

	r.ready <<= 8
	r.offset += 8
	r.available -= 8

	if r.available <= 0 {
		return result
	}

	if r.offset == 32 {
		r.offset = 0
		r.index++
		r.ready = r.buff[r.index]
		return result
	}

	if r.offset > 32 {
		r.offset -= 32
		r.index++
		next := r.buff[r.index]
		r.ready = next << r.offset
		return result | byte(next>>(32-r.offset))
	}

	return result
}

func (r *reader) ReadByte() (byte, error) {
	if r.available <= 0 {
		return 0x00, io.EOF
	}

	switch r.direction {
	case Up:
		return r.readUp(), nil
	case Down:
		return r.readDown(), nil
	case Left:
		return r.readLeft(), nil
	case Right:
		return r.readRight(), nil
	default:
		panic("Unrecognised read direction")
	}
}

func (m *Matrix) Reader(row, col int, direction Direction, count int) io.ByteReader {
	if row < 0 || row >= m.height || col < 0 || col >= m.width {
		panic("Arg out of range")
	}

	index := (row * m.intsPerRow) + (col / 32)
	offset := col % 32

	result := &reader{
		buff:   m.bits,
		index:  index,
		offset: offset,
	}

	switch direction {
	case Right:
		result.available = m.width - col
		result.ready = m.bits[index] << offset
		result.direction = Right
	case Left:
		result.available = col + 1
		result.ready = m.bits[index] >> (31 - offset)
		result.direction = Left
	case Up:
		result.available = row + 1
		result.bit = 0x80000000 >> offset
		result.direction = Up
	case Down:
		result.available = m.height - row
		result.bit = 0x80000000 >> offset
		result.direction = Down
	default:
		panic("Unrecognised reader direction")
	}

	if result.available > count {
		result.available = count
	}

	return result
}
