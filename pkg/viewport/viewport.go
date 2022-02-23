package viewport

import (
	"github.com/realency/arke/pkg/display"
)

// ViewPort is the interface for types that connect a Canvas to a physical dot-matrix display.
//
// A ViewPort can be thought of as a frame surrounding a subregion of pixels in a Canvas.
// Pixels within the frame are represented in the display, and changes made to the canvas effect
// the same change in the display.  Allowing a viewport to attach to a subregion of the canvas
// enables effects such as scrolling, in which the viewport moves over the underlying canvas.
type ViewPort interface {
	// Attach attaches the ViewPort to a canvas at a specific location, so that changes in region framed by the canvas are reflected in the display.
	Attach(canvas *display.Canvas, row, col int)

	// Detach detaches the ViewPort from the canvas it is currently attached to.
	Detach()

	// Locate repositions the ViewPort at a new position on the underlying canvas.
	Locate(row, col int)

	// Offset returns the current location of the ViewPort - its offset into the canvas.
	Offset() (row, col int)

	// Size returns the size of the ViewPort in pixels.
	Size() (height, width int)

	// Canvas returns the attached canvas, or nil if the ViewPort is not attached to any canvas.
	Canvas() *display.Canvas
}
