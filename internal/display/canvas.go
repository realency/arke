package display

import . "github.com/realency/arke/pkg/display"

type canvas struct {
	bits      [][]bool
	observers []CanvasObserver
	batching  bool
}

func NewCanvas(height, width int) Canvas {
	result := &canvas{
		bits:     make([][]bool, height),
		batching: false,
	}

	for i := range result.bits {
		result.bits[i] = make([]bool, width)
	}

	return result
}

func (c *canvas) Get(row, col int) bool {
	return c.bits[row][col]
}

func (c *canvas) Height() int {
	return len(c.bits)
}

func (c *canvas) Width() int {
	return len(c.bits[0])
}

func (c *canvas) Set(row, col int, value bool) {
	c.bits[row][col] = value
	if !c.batching {
		c.notify()
	}
}

func (c *canvas) Write(from [][]bool, row, col int) {
	for i, r := range from {
		if i+row >= len(c.bits) {
			break
		}

		for j, b := range r {
			if j+col >= len(c.bits[0]) {
				break
			}

			c.bits[i+row][j+col] = b
		}
	}

	if !c.batching {
		c.notify()
	}
}

func (c *canvas) Observe(observer CanvasObserver) {
	c.observers = append(c.observers, observer)
}

func (c *canvas) StartUpdate() {
	c.batching = true
}

func (c *canvas) EndUpdate() {
	c.batching = false
	c.notify()
}

func (c *canvas) Clear() {
	for i, r := range c.bits {
		c.bits[i] = make([]bool, len(r))
	}

	if !c.batching {
		c.notify()
	}
}

func (c *canvas) notify() {
	for _, o := range c.observers {
		o(c.bits)
	}
}
