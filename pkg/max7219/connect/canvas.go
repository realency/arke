package connect

import (
	internal "github.com/realency/arke/internal/display"

	"github.com/realency/arke/pkg/display"
)

func NewCanvas(height, width int) display.Canvas {
	return internal.NewCanvas(height, width)
}
