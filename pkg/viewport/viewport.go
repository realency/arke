package viewport

import (
	"github.com/realency/arke/pkg/display"
)

type ViewPort interface {
	Attach(canvas *display.Canvas, row, col int)
	Detach()
	Locate(row, col int)
	Offset() (row, col int)
	Size() (height, width int)
	Canvas() *display.Canvas
}
