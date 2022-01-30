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
	c.canvas.BeginUpdate()
	defer c.canvas.EndUpdate()

	i := 0
	for _, r := range asStr {
		if c.col >= c.canvas.Width() {
			break
		}
		m := c.font(r)
		width := m.Width()
		c.canvas.Write(m, c.row, c.col)
		c.col += width
		i++
	}

	return i, nil
}
