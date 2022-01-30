package bits

import "strings"

type Matrix struct {
	bits       []uint32
	height     int
	width      int
	intsPerRow int
}

func NewMatrix(height, width int) *Matrix {
	var rowLen int = ((width - 1) / 32) + 1
	return &Matrix{
		bits:       make([]uint32, rowLen*height),
		height:     height,
		width:      width,
		intsPerRow: rowLen,
	}
}

func (m *Matrix) Height() int {
	return m.height
}

func (m *Matrix) Width() int {
	return m.width
}

func (m *Matrix) selector(row, col int) (index int, mask uint32) {
	if row < 0 || row >= m.height || col < 0 || col >= m.width {
		panic("Arg out fo range")
	}
	return (m.intsPerRow * row) + (col / 32), uint32(0x80000000) >> (col % 32)
}

func (m *Matrix) Get(row, col int) bool {
	i, k := m.selector(row, col)
	return (m.bits[i] & k) != 0
}

func (m *Matrix) Set(row, col int, value bool) {
	i, k := m.selector(row, col)
	if value {
		m.bits[i] |= k
	} else {
		m.bits[i] &= (k ^ 0xFFFFFFFF)
	}
}

func (m *Matrix) Clear() {
	m.bits = make([]uint32, m.intsPerRow*m.height)
}

func (m *Matrix) Clone() *Matrix {
	result := &Matrix{
		bits:       make([]uint32, len(m.bits)),
		height:     m.height,
		width:      m.width,
		intsPerRow: m.intsPerRow,
	}
	copy(result.bits, m.bits)
	return result
}

func Copy(from *Matrix, fromRow, fromCol int, to *Matrix, toRow, toCol, height, width int) (int, int) {
	if height == 0 || width == 0 {
		return 0, 0
	}

	if height < 0 || width < 0 {
		panic("Arg out of bounds")
	}

	if fromRow < 0 || fromRow >= from.height || fromCol < 0 || fromCol >= from.width {
		panic("Arg out of bounds")
	}

	if toRow < 0 || toRow >= to.height || toCol < 0 || toCol >= to.width {
		panic("Arg out of bounds")
	}

	if height > from.height-fromRow {
		height = from.height - fromRow
	}

	if height > to.height-toRow {
		height = to.height - toRow
	}

	if width > from.width-fromCol {
		width = from.width - fromCol
	}

	if width > to.width-toCol {
		width = to.width - toCol
	}

	for i := 0; i < height; i++ {
		r := from.Reader(fromRow+i, fromCol, Right, width)
		w := to.Writer(toRow+i, toCol, width)
		for {
			b, e := r.ReadByte()
			if e != nil {
				break
			}
			w.WriteByte(b)
			if e != nil {
				break
			}
		}
	}

	return height, width
}

func (m *Matrix) String() string {
	var sb strings.Builder
	for i := 0; i < m.height; i++ {
		for j := 0; j < m.width; j++ {
			if m.Get(i, j) {
				sb.WriteString("@ ")
			} else {
				sb.WriteString(". ")
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
