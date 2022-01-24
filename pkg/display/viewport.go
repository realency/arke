package display

type ViewPort interface {
	Attach(canvas *Canvas, row, col int)
}
