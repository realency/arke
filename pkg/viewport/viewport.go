package viewport

import (
	"errors"

	"github.com/realency/arke/pkg/display"
)

type ViewPort interface {
	Attach(canvas *display.Canvas, row, col int) error
	Detach() error
	Shift(rowDelta, colDelta int) error
	Offset() (row, col int)
	Size() (height, width int)
	Canvas() *display.Canvas
}

var (
	EAlreadyAttached = errors.New("ViewPort already attached to a canvas")
	ENotAttached     = errors.New("ViewPort not attached to a canvas")
)
