package bits_test

import (
	"testing"

	"github.com/realency/arke/pkg/bits"
)

type coord struct {
	row, col int
}

type size struct {
	height, width int
}

func initMatrix(s size, ones []coord) *bits.Matrix {
	result := bits.NewMatrix(s.height, s.width)
	for _, c := range ones {
		result.Set(c.row, c.col, true)
	}
	return result
}

func TestNewMatrixCreatesMatrixWithExpectedDimensions(t *testing.T) {
	sizes := []size{
		{2, 2},
		{2, 10},
		{30000, 30000},
		{1000000, 2},
		{2, 1000000},
	}

	for _, s := range sizes {
		h, w := bits.NewMatrix(s.height, s.width).Size()
		if s.height != h || s.width != w {
			t.Error("Reported size of matrix not equal to requested size")
		}
	}
}

func TestNewMatrixPanicsForArgumentsOutOfBounds(t *testing.T) {
	caller := func(h, w int) {
		defer func() {
			recover()
		}()
		bits.NewMatrix(h, w)
		t.Errorf("Did not panic for %d, %d", h, w)
	}

	sizes := []size{
		{-7, 10},
		{-8, -5},
		{9, -32},
		{0, -9},
		{-1, 0},
	}

	for _, s := range sizes {
		caller(s.height, s.width)
	}
}

func TestNewMatrixReturnsSameMatrixForAnyZeroSize(t *testing.T) {
	sizes := []size{
		{0, 0},
		{0, 32},
		{67, 0},
	}

	for _, s := range sizes {
		if bits.NewMatrix(s.height, s.width) != bits.ZeroMatrix {
			t.Errorf("NewMatrix did not return ZeroMatrix for %d, %d", s.height, s.width)
		}
	}
}

func TestNewMatrixReturnsMatrixWithOnlyZeroBits(t *testing.T) {
	m := bits.NewMatrix(100, 100)

	for row := 0; row < 100; row++ {
		for col := 0; col < 100; col++ {
			if m.Get(row, col) {
				t.Error("New matrix includes at least one one-bit")
			}
		}
	}
}

func TestGetReturnsValuesSetBySet(t *testing.T) {
	m := bits.NewMatrix(20, 20)

	data := []struct {
		row, col int
		value    bool
	}{
		{0, 0, true},
		{2, 5, false},
		{10, 0, true},
	}

	lookup := func(row, col int) bool {
		for _, d := range data {
			if row == d.row && col == d.col {
				return d.value
			}
		}
		return false
	}

	for _, d := range data {
		m.Set(d.row, d.col, d.value)
	}

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			expected := lookup(i, j)
			actual := m.Get(i, j)
			if expected != actual {
				t.Errorf("Value m[%d,%d] was %v, when %v was expected.", i, j, actual, expected)
			}
		}
	}
}

func TestAllBitsAreZeroAfterClear(t *testing.T) {
	m := initMatrix(size{32, 100}, []coord{
		{10, 12},
		{25, 0},
		{0, 99},
		{7, 64},
		{6, 63},
	})

	m.Clear()

	for i := 0; i < 32; i++ {
		for j := 0; j < 100; j++ {
			if m.Get(i, j) {
				t.Errorf("1-bit at [%d,%d] after Clear()", i, j)
			}
		}
	}
}

func TestCloneCreatesIdenticalCopyOfMatrix(t *testing.T) {
	ones := []coord{
		{0, 0},
		{0, 1},
		{90, 18},
		{5, 31},
		{10, 10},
	}
	m := initMatrix(size{100, 32}, ones)
	clone := m.Clone()

	mh, mw := m.Size()
	ch, cw := clone.Size()

	if mh != ch || mw != cw {
		t.Error("Source and clone matrices do not have same size")
	}

	for i := 0; i < mh; i++ {
		for j := 0; j < mw; j++ {
			if m.Get(i, j) != clone.Get(i, j) {
				t.Error("Source and clone matrices do not have the same state")
			}
		}
	}
}

func TestCloneIsNotReferenceIdenticalToSource(t *testing.T) {
	m := bits.NewMatrix(20, 20)
	clone := m.Clone()
	m.Set(8, 9, true)
	if clone.Get(8, 9) {
		t.Error("Clone matrix affected by change to source matrix")
	}
}
