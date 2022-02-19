package viewport

import (
	"github.com/realency/arke/pkg/display"
)

// ViewPort is the interface for types that connect a Canvas to a physical dot-matrix display.
//
// A ViewPort can be conceptualised as a frame surrounding a subregion of pixels in a Canvas.
// Pixels within the frame are represented in the display, and changes made to the canvas effect
// the same change in the display.  Allowing a viewport to attach to a subregion of the canvas
// allows for effects such as scrolling, in which the viewport moves over the underlying canvas.
type ViewPort interface {
	Attach(canvas *display.Canvas, row, col int)
	Detach()
	Locate(row, col int)
	Offset() (row, col int)
	Size() (height, width int)
	Canvas() *display.Canvas
}
