package viewport

import "github.com/realency/arke/pkg/display"

type ViewPort interface {
	Attach(canvas *display.Canvas, row, col int)
	Shift(rowDelta, colDelta int)
}
