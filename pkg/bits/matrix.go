package bits

import (
	"strings"
)

// A Matrix is a two-dimensional array of bits with fixed height and width.
//
// Matrix is not thread-safe.  Read-only operations such as Get and Clone
// may safely be called concurrently.  Write operations such as Set and Clear
// should not be called concurrently with each other, or with read operations.
type Matrix struct {
	bits          []uint32 // And array of unsigned ints as a bitfield storing the individual bits of the matrix
	height, width int
	intsPerRow    int // Number of elements of the bits slice per row of the matrix, calculated when the Matrix is initialised and stored for reuse
}

// NewMatrix creates a new Matrix with a give size
func NewMatrix(height, width int) *Matrix {
	var rowLen int = ((width - 1) / 32) + 1
	return &Matrix{
		bits:       make([]uint32, rowLen*height),
		height:     height,
		width:      width,
		intsPerRow: rowLen,
	}
}

// Size returns the size of the Matrix as a two-tuple of height and width
func (m *Matrix) Size() (height, width int) {
	return m.height, m.width
}

// Get returns the state of a specific bit in the Matrix.
// Arguments specify the coordinates of the bit.  Get will panic if
// the arguments are out of bounds.
func (m *Matrix) Get(row, col int) bool {
	i, k := m.selector(row, col)
	return (m.bits[i] & k) != 0
}

// Set allocates state to a specific bit in the Matrix.
// Arguments specify the coordinates of the bit and the required state.
// Set will panic if the arguments are out of bounds.
func (m *Matrix) Set(row, col int, value bool) {
	i, k := m.selector(row, col)
	if value {
		m.bits[i] |= k
	} else {
		m.bits[i] &= (k ^ 0xFFFFFFFF)
	}
}

// Clear resets all the bits in the matrix back to zero.
func (m *Matrix) Clear() {
	m.bits = make([]uint32, m.intsPerRow*m.height)
}

// Clone creates an exact copy of the Matrix in its current state.
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

// Copy copies a sub-range of bits from one matrix to another.  All the bits in the range are overwritten
// in the destination matrix.
//
// Copies from source matrix, at origin (sourceRow, sourceCol) to the dest matrix ar (destRow, destCol).
// Copy will panic if either the source origin or the destination origin are out of bounds.
// Copies a rectangle up to the size given by height and width.  If the maximum-sized rectange exceeds
// the bonds of either the source or destination matrix, it is trimmed.
// Returns the actual height and width of the rectange copied as a result.
//
// Copying is a read-only operation with respect to the source matrix, with the usual implications for
// concurrency.
func Copy(source *Matrix, sourceRow, sourceCol int, dest *Matrix, destRow, destCol, height, width int) (int, int) {
	// A whole load of bounds-checking and trimming logic
	if height == 0 || width == 0 {
		return 0, 0
	}
	if height < 0 || width < 0 {
		panic("Arg out of bounds")
	}
	if sourceRow < 0 || sourceRow >= source.height || sourceCol < 0 || sourceCol >= source.width {
		panic("Arg out of bounds")
	}
	if destRow < 0 || destRow >= dest.height || destCol < 0 || destCol >= dest.width {
		panic("Arg out of bounds")
	}
	if height > source.height-sourceRow {
		height = source.height - sourceRow
	}
	if height > dest.height-destRow {
		height = dest.height - destRow
	}
	if width > source.width-sourceCol {
		width = source.width - sourceCol
	}
	if width > dest.width-destCol {
		width = dest.width - destCol
	}

	// The acual copy operation is performed using stream readers on the source and writers on the destination
	for i := 0; i < height; i++ {
		r := newStream(source, sourceRow+i, sourceCol, width)
		w := newStream(source, destRow+i, destCol, width)
		streamCopy(r, w)
	}

	return height, width
}

// String generates a string representation of the matrix.
// The string generated is intended for visual inspection, and applies
// a style appropriate to that.
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

// A utility function to find the index into the m.bits array and the appropriate
// Bitwise mask to select the bit, given a bit's coordinates.
// Panics if the arguments are out of range.
func (m *Matrix) selector(row, col int) (index int, mask uint32) {
	if row < 0 || row >= m.height || col < 0 || col >= m.width {
		panic("Arg out fo range")
	}
	return (m.intsPerRow * row) + (col / 32), uint32(0x80000000) >> (col % 32)
}
