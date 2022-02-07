package bits

import "fmt"

type stream struct {
	buff          []uint32
	index, offset int
	available     int
}

func newStream(m *Matrix, row, col, count int) *stream {
	if row < 0 || row >= m.height || col < 0 || col >= m.width || col+count > m.width {
		panic(fmt.Sprintf("Arguments out of bounds.  row=%d, col=%d, count=%d, m.height=%d, m.width=%d", row, col, count, m.height, m.width))
	}

	return &stream{
		buff:      m.bits,
		index:     (m.intsPerRow * row) + (col / 32),
		offset:    col % 32,
		available: count,
	}
}

func (r *stream) usable(count int) int {
	n := count
	if n > 32 {
		n = 32
	}
	if r.available < n {
		n = r.available
	}
	if (32 - r.offset) < n {
		n = 32 - r.offset
	}
	return n
}

func (r *stream) seek(count int) {
	r.available -= count
	r.offset += count
	if r.offset == 32 {
		r.offset = 0
		r.index++
	}
}

func (r *stream) read(dest *uint32, count int) int {
	count = r.usable(count)

	switch count {
	case 0:
		return 0
	case 32:
		// Optimisation for common case where we're reading the whole int.
		*dest = r.buff[r.index]
		r.index++
		r.available -= 32
		return 32
	default:
	}

	*dest = r.buff[r.index] << r.offset
	r.seek(count)
	return count
}

func (r *stream) write(source uint32, count int) int {
	count = r.usable(count)
	switch count {
	case 0:
		return 0
	case 32:
		// Optimisation for the common case where we're writing the whole int
		r.buff[r.index] = source
		r.available -= 32
		r.index++
		return 32
	default:
	}

	var mask uint32 = (0xFFFFFFFF << (32 - count)) >> r.offset
	source = (source >> r.offset) & mask
	mask ^= 0xFFFFFFFF
	mask &= r.buff[r.index]
	r.buff[r.index] = mask | source

	r.seek(count)
	return count
}

func streamCopy(source, dest *stream) int {
	var buff uint32
	count := 0
	for {
		var readCount int
		if readCount = source.read(&buff, 32); readCount == 0 {
			return count
		}

		for readCount > 0 {
			var writeCount int
			if writeCount = dest.write(buff, readCount); writeCount == 0 {
				return count
			}
			count += writeCount
			readCount -= writeCount
		}
	}
}
