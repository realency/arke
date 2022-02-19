package bits_test

import (
	"testing"

	"github.com/realency/arke/pkg/bits"
)

func TestNewMatrixCreatesMatrixWithExpectedDimensions(t *testing.T) {
	sizes := []struct {
		height, width int
	}{
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

	sizes := []struct {
		height, width int
	}{
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
	sizes := []struct {
		height, width int
	}{
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
