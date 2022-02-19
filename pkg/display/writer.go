package display

import "io"

type canvasWriter struct {
	canvas *Canvas
	font   Font
	row    int
	col    int
}

// NewWriter returns a new io.Writer for writing text to a Canvas.
// The writer writes unicode text in a left-to-right direction to the Canvas from the given location using a given font.
func NewWriter(canvas *Canvas, font Font, row, col int) io.Writer {
	return &canvasWriter{
		canvas: canvas,
		font:   font,
		row:    row,
		col:    col,
	}
}

func (c *canvasWriter) Write(p []byte) (n int, err error) {
	asStr := string(p)
	_, w := c.canvas.Size()

	c.canvas.BeginUpdate()
	defer c.canvas.EndUpdate()

	i := 0
	for _, r := range asStr {
		if c.col >= w {
			break
		}
		m := c.font(r)
		_, width := m.Size()
		c.canvas.Write(m, c.row, c.col)
		c.col += width
		i++
	}

	return i, nil
}
