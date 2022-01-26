package display

type canvasWriter struct {
	canvas *Canvas
	font   Font
	row    int
	col    int
}

func NewWriter(canvas *Canvas, font Font, row, col int) *canvasWriter {
	return &canvasWriter{
		canvas: canvas,
		font:   font,
		row:    row,
		col:    col,
	}
}

func (c *canvasWriter) Write(p []byte) (n int, err error) {
	asStr := string(p)
	c.canvas.StartUpdate()
	defer c.canvas.EndUpdate()

	i := 0
	for _, r := range asStr {
		bits := c.font(r)
		width := len(bits[0])
		if c.canvas.Width() > c.col+width {
			break
		}
		c.canvas.Write(bits, c.row, c.col)
		c.col += width
		i++
	}

	return i, nil
}
