package display

import "github.com/realency/arke/internal/bits"

type CanvasObserver chan<- *bits.ImmutableBuffer

type Canvas struct {
	buff        *bits.Buffer
	observers   []CanvasObserver
	updateLevel uint
}

func NewCanvas(height, width int) *Canvas {
	return &Canvas{
		buff:        bits.NewBuffer(height, width),
		updateLevel: 0,
	}
}

func (c *Canvas) Get(row, col int) bool {
	return c.buff.Get(row, col)
}

func (c *Canvas) Height() int {
	return c.buff.Height()
}

func (c *Canvas) Width() int {
	return c.buff.Width()
}

func (c *Canvas) Set(row, col int, value bool) {
	c.buff.Set(row, col, value)
}

func (c *Canvas) Write(from [][]bool, row, col int) {
	h := c.Height()
	w := c.Width()
	for i, r := range from {
		if i+row >= h {
			break
		}

		for j, b := range r {
			if j+col >= w {
				break
			}

			c.buff.Set(i, j, b)
		}
	}

	if c.updateLevel == 0 {
		c.notify()
	}
}

func (c *Canvas) Observe(observer CanvasObserver) {
	c.observers = append(c.observers, observer)
}

func (c *Canvas) StartUpdate() {
	c.updateLevel++
}

func (c *Canvas) EndUpdate() {
	if c.updateLevel == 0 {
		panic("EndUpdate called out of sequence")
	}
	c.updateLevel--
	if c.updateLevel == 0 {
		c.notify()
	}
}

func (c *Canvas) notify() {
	for _, o := range c.observers {
		select {
		case o <- c.buff.GetImmutableCopy():
		default:
		}
	}
}
