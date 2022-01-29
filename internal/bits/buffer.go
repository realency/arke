package bits

import (
	"io"
	"sync"
)

type Buffer struct {
	height      int
	width       int
	bits        []byte
	dirty       []byte
	bytesPerRow int
	mutex       *sync.Mutex
}

type ImmutableBuffer struct {
	height      int
	width       int
	bits        []byte
	dirty       []byte
	bytesPerRow int
}

func NewBuffer(height, width int) *Buffer {
	if width <= 0 || height <= 0 {
		panic("arg out of bounds")
	}
	rowLen := ((width - 1) / 8) + 1

	return &Buffer{
		height:      height,
		width:       width,
		bits:        make([]byte, rowLen*height),
		dirty:       make([]byte, rowLen*height),
		bytesPerRow: rowLen,
		mutex:       &sync.Mutex{},
	}
}

func (m *Buffer) selector(row, col int) (int, byte) {
	return (row * int(m.bytesPerRow)) + (col / 8), 0x80 >> (col % 8)
}

func (m *Buffer) Get(row, col int) bool {
	idx, mask := m.selector(row, col)
	return (m.bits[idx] & mask) != 0
}

func (m *Buffer) Set(row, col int, value bool) {
	idx, mask := m.selector(row, col)
	m.mutex.Lock()
	if (m.bits[idx]&mask != 0x00) != value {
		m.bits[idx] ^= mask
		m.dirty[idx] ^= mask
	}
	m.mutex.Unlock()
}

func (m *Buffer) Flip(row, col int) {
	idx, mask := m.selector(row, col)
	m.mutex.Lock()
	m.bits[idx] ^= mask
	m.dirty[idx] ^= mask
	m.mutex.Unlock()
}

func (m *Buffer) ResetDirtyFlags() {
	m.mutex.Lock()
	m.dirty = make([]byte, m.bytesPerRow*m.height)
	m.mutex.Unlock()
}

func (m *Buffer) GetImmutableCopy() *ImmutableBuffer {
	result := &ImmutableBuffer{
		height:      m.height,
		width:       m.width,
		bits:        make([]byte, m.bytesPerRow*m.height),
		dirty:       make([]byte, m.bytesPerRow*m.height),
		bytesPerRow: m.bytesPerRow,
	}
	m.mutex.Lock()
	copy(result.bits, m.bits)
	copy(result.dirty, m.dirty)
	m.mutex.Unlock()
	return result
}

func (m *Buffer) RowReader(row, col int) io.Reader {
	if col > m.width || row < 0 || col < 0 || row > m.height {
		panic("arg out of bounds")
	}
	a := make([]byte, m.bytesPerRow)
	m.mutex.Lock()
	copy(a, m.bits[row*m.bytesPerRow:])
	m.mutex.Unlock()
	return NewArrayReader(a, col)
}
