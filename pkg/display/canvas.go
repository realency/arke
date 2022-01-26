package display

type CanvasObserver chan<- [][]bool

type Canvas struct {
	bits      [][]bool
	observers []CanvasObserver
	batching  bool
}

func NewCanvas(height, width int) *Canvas {
	result := &Canvas{
		bits:     make([][]bool, height),
		batching: false,
	}

	for i := range result.bits {
		result.bits[i] = make([]bool, width)
	}

	return result
}

func (c *Canvas) Get(row, col int) bool {
	return c.bits[row][col]
}

func (c *Canvas) Height() int {
	return len(c.bits)
}

func (c *Canvas) Width() int {
	return len(c.bits[0])
}

func (c *Canvas) Set(row, col int, value bool) {
	c.bits[row][col] = value
	if !c.batching {
		c.notify()
	}
}

func (c *Canvas) Write(from [][]bool, row, col int) {
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

func (c *Canvas) Observe(observer CanvasObserver) {
	c.observers = append(c.observers, observer)
}

func (c *Canvas) StartUpdate() {
	c.batching = true
}

func (c *Canvas) EndUpdate() {
	c.batching = false
	c.notify()
}

func (c *Canvas) Clear() {
	for i, r := range c.bits {
		c.bits[i] = make([]bool, len(r))
	}

	if !c.batching {
		c.notify()
	}
}

func (c *Canvas) notify() {
	for _, o := range c.observers {
		select {
		case o <- c.bits: // TODO: Danger!  This is by-reference.  Bits could change by the time they're read.  Address this
		default:
		}
	}
}
