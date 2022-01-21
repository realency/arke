package display

type CanvasObserver func([][]bool)

type Canvas interface {
	Height() int
	Width() int
	Get(row, col int) bool
	Set(row, col int, value bool)
	Write(from [][]bool, row, col int)
	Observe(observer CanvasObserver)
	StartUpdate()
	EndUpdate()
	Clear()
}
