package display

import (
	"sync"
	"sync/atomic"

	"github.com/realency/arke/pkg/bits"
)

type CanvasObserver chan<- CanvasUpdate

type Canvas struct {
	buff      *bits.Matrix
	observers map[uint64]CanvasObserver
	mutex     *sync.RWMutex
	update    CanvasUpdateKind
	nextId    uint64
}

type CanvasUpdateKind byte

const (
	CanvasNoOp  CanvasUpdateKind = 0x00
	CanvasClear CanvasUpdateKind = 0x01
	CanvasWrite CanvasUpdateKind = 0x02
	CanvasBatch CanvasUpdateKind = 0x04
)

type CanvasUpdate struct {
	Buff *bits.Matrix
	Kind CanvasUpdateKind
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
	defer c.mutex.Unlock()
	c.buff.Set(row, col, value)
	c.updated(CanvasWrite)
}

func (c *Canvas) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.buff.Clear()
	c.updated(CanvasClear)
}

func (c *Canvas) Matrix() *bits.Matrix {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.buff.Clone()
}

func (c *Canvas) Write(source *bits.Matrix, row, col int) {
	h, w := source.Size()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	bits.Copy(source, 0, 0, c.buff, row, col, h, w)
	c.updated(CanvasWrite)
}

func (c *Canvas) AddObserver(observer CanvasObserver) uint64 {
	i := atomic.AddUint64(&c.nextId, 1)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.observers[i] = observer
	return i
}

func (c *Canvas) RemoveObserver(id uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.observers, id)
}

func (c *Canvas) BeginUpdate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.update != CanvasNoOp {
		panic("BeginUpdate called out of sequence - update already underway")
	}
	c.update = CanvasBatch
}

func (c *Canvas) EndUpdate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.update == CanvasNoOp {
		panic("EndUpdate called out of sequence - no update underway")
	}
	c.notify(c.update)
	c.update = CanvasNoOp
}

func (c *Canvas) updated(kind CanvasUpdateKind) {
	if (c.update & CanvasBatch) != 0 {
		c.update |= kind
	} else {
		c.notify(kind)
	}
}

func (c *Canvas) notify(kind CanvasUpdateKind) {
	clone := c.buff.Clone()
	for _, o := range c.observers {
		select {
		case o <- CanvasUpdate{clone, kind}:
		default:
		}
	}
}
