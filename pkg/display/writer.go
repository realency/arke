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
