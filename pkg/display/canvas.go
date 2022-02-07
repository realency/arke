package display

import (
	"sync"

	"github.com/realency/arke/pkg/bits"
)

type CanvasObserver chan<- *bits.Matrix

type Canvas struct {
	buff        *bits.Matrix
	observers   map[uint64]CanvasObserver
	mutex       *sync.RWMutex
	updateDepth uint32
	observerId  uint64
}

func NewCanvas(height, width int) *Canvas {
	return &Canvas{
		buff:      bits.NewMatrix(height, width),
		observers: make(map[uint64]CanvasObserver),
		mutex:     &sync.RWMutex{},
	}
}

func (c *Canvas) Size() (height, width int) {
	return c.buff.Size()
}

func (c *Canvas) Get(row, col int) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.buff.Get(row, col)
}

func (c *Canvas) Set(row, col int, value bool) {
	c.mutex.Lock()
	defer func() {
		c.updated()
		c.mutex.Unlock()
	}()
	c.buff.Set(row, col, value)
}

func (c *Canvas) Clear() {
	c.mutex.Lock()
	defer func() {
		c.updated()
		c.mutex.Unlock()
	}()
	c.buff.Clear()
}

func (c *Canvas) Matrix() *bits.Matrix {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.buff.Clone()
}

func (c *Canvas) Write(source *bits.Matrix, row, col int) {
	h, w := source.Size()
	c.mutex.Lock()
	defer func() {
		c.updated()
		c.mutex.Unlock()
	}()
	bits.Copy(source, 0, 0, c.buff, row, col, h, w)
}

func (c *Canvas) AddObserver(observer CanvasObserver) (uint64, *bits.Matrix) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.observerId++
	i := c.observerId
	c.observers[i] = observer
	return i, c.buff.Clone()
}

func (c *Canvas) RemoveObserver(id uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.observers, id)
}

func (c *Canvas) BeginUpdate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.updateDepth++
}

func (c *Canvas) EndUpdate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.updateDepth == 0 {
		panic("EndUpdate called out of sequence")
	}
	c.updateDepth--
	if c.updateDepth == 0 {
		c.updated()
	}
}

func (c *Canvas) updated() {
	if c.updateDepth > 0 {
		return
	}

	b := c.buff.Clone()

	for _, o := range c.observers {
		select {
		case o <- b:
		default:
		}
	}
}
