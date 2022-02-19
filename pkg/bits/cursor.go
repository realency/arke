package bits

// A Cursor allows a BitMatrix to be traversed by reading across, or up and down the matrix.
//
// The type provides a mechanism to traverse a bit matrix that is optimised for reading contiguous bits,
// as compared to repeatedly using Get() in a loop.
//
// When reading from left to right, or downwards, the read is sequenced as follows: sample current bit then advance cursor position.
// When reading from right to left, or upwards, the read is sequenced as follows: advance cursor position then sample bit at new position.
// The different sequencing allows the expected behaviour that reading left then right repeatedly keeps returning the same bit.
type Cursor struct {
	matrix   *Matrix
	row, col int
	index    int
	mask     uint32
	current  uint32
}

// NewCursor returns a new instance of a Cursor.
//
// Receives the matrix to address as an argument, as well as the starting location.
// For rightward or downward reads, the initial cursor position should be on the first bit to be read.
// For leftward or upward reads, the initial cursor position should be one before the bit to be read (to its right or beneath it).
// Because of this, positioning the cursor at row = height, or at col = width is permitted.
func NewCursor(m *Matrix, row, col int) *Cursor {
	if row < 0 || row > m.height || col < 0 || col > m.width {
		panic("Arg out of bounds")
	}
	i := (row * m.intsPerRow) + (col / 32)
	return &Cursor{
		matrix:  m,
		row:     row,
		col:     col,
		index:   i,
		mask:    0x80000000 >> (col % 32),
		current: m.bits[i],
	}
}

// Position returns the current position of the cursor.
func (c *Cursor) Position() (row, col int) {
	return c.row, c.col
}

// ReadLeft positions the cursor one bit to the left and returns the value at the new position.
// if ok is returned as false, the cursor cannot move further left, it's already on column zero.
func (c *Cursor) ReadLeft() (bit, ok bool) {
	if c.col == 0 || c.row == c.matrix.height {
		return false, false
	}
	c.col--
	if c.mask <<= 1; c.mask == 0 {
		c.index--
		c.current = c.matrix.bits[c.index]
		c.mask = 0x00000001
	}
	return c.current&c.mask != 0, true
}

// ReadRight samples the bit at the current position and then positions the cursor one bit to the right.
// if ok is returned as false, the cursor cannot move further right, it's already past the rightmost column.
func (c *Cursor) ReadRight() (bit, ok bool) {
	if c.col == c.matrix.width || c.row == c.matrix.height {
		return false, false
	}
	result := c.current&c.mask != 0

	c.col++
	if c.mask >>= 1; c.mask == 0 {
		c.index++
		c.current = c.matrix.bits[c.index]
		c.mask = 0x80000000
	}

	return result, true
}

// ReadUp positions the cursor one bit further up and returns the value at the new position.
// if ok is returned as false, the cursor cannot move further up, it's already on row zero.
func (c *Cursor) ReadUp() (bit, ok bool) {
	if c.row == 0 || c.col == c.matrix.width {
		return false, false
	}

	c.row--
	c.index -= c.matrix.intsPerRow
	c.current = c.matrix.bits[c.index]
	return c.current&c.mask != 0, true
}

// ReadDown samples the bit at the current position and then positions the cursor one bit further down.
// if ok is returned as false, the cursor cannot move further down, it's already past the bottom row.
func (c *Cursor) ReadDown() (bit, ok bool) {
	if c.row == c.matrix.height || c.col == c.matrix.width {
		return false, false
	}

	result := c.current&c.mask != 0

	c.row++
	c.index += c.matrix.intsPerRow
	c.current = c.matrix.bits[c.index]

	return result, true
}

// ReadLeftByte constructs a byte from 8 consecutive reads leftward.
// Because of the leftward direction, bits in the result are sequenced in the reverse order of the natural order in the matrix.
// That is to say that the most significant bit of the returned byte is read from the rightmost position read by the cursor.
// Returns the resulting byte and the number of bits successfully read before reading past column zero.
// The last bit read is always in the least significant position, regardless of the number of bits returned.
func (c *Cursor) ReadLeftByte() (b byte, bits int) {
	bits = 0
	b = 0x00
	for bits < 8 {
		bit, ok := c.ReadLeft()
		if !ok {
			break
		}
		b <<= 1
		if bit {
			b |= 0x01
		}
		bits++
	}
	return
}

// ReadRightByte constructs a byte from 8 consecutive reads rightward.
// Returns the resulting byte and the number of bits successfully read before reading past the width of the matrix.
// The last bit read is always in the least significant position, regardless of the number of bits returned.
func (c *Cursor) ReadRightByte() (b byte, bits int) {
	bits = 0
	b = 0x00
	for bits < 8 {
		bit, ok := c.ReadRight()
		if !ok {
			break
		}
		b <<= 1
		if bit {
			b |= 0x01
		}
		bits++
	}
	return
}

// ReadUpByte constructs a byte from 8 consecutive reads upward.
// Returns the resulting byte and the number of bits successfully read before reading past row zero.
// The last bit read is always in the least significant position, regardless of the number of bits returned.
func (c *Cursor) ReadUpByte() (b byte, bits int) {
	bits = 0
	b = 0x00
	for bits < 8 {
		bit, ok := c.ReadUp()
		if !ok {
			break
		}
		b <<= 1
		if bit {
			b |= 0x01
		}
		bits++
	}
	return
}

// ReadDownByte constructs a byte from 8 consecutive reads downward.
// Returns the resulting byte and the number of bits successfully read before reading past the height of the matrix.
// The last bit read is always in the least significant position, regardless of the number of bits returned.
func (c *Cursor) ReadDownByte() (b byte, bits int) {
	bits = 0
	b = 0x00
	for bits < 8 {
		bit, ok := c.ReadDown()
		if !ok {
			break
		}
		b <<= 1
		if bit {
			b |= 0x01
		}
		bits++
	}
	return
}
