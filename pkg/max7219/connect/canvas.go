package connect

import (
	internal "realency/arke/internal/display"

	"realency/arke/pkg/display"
)

func NewCanvas(height, width int) display.Canvas {
	return internal.NewCanvas(height, width)
}
