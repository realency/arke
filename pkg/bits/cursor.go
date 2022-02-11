package bits

type Cursor struct {
	matrix   *Matrix
	row, col int
	index    int
	mask     uint32
	current  uint32
}

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

func (c *Cursor) Position() (row, col int) {
	return c.row, c.col
}

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

func (c *Cursor) ReadUp() (bit, ok bool) {
	if c.row == 0 || c.col == c.matrix.width {
		return false, false
	}

	c.row--
	c.index -= c.matrix.intsPerRow
	c.current = c.matrix.bits[c.index]
	return c.current&c.mask != 0, true
}

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
