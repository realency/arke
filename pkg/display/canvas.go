package display

import (
	"sync"

	"github.com/realency/arke/pkg/bits"
)

// CanvasObserver is an alias for a channel on which to receive notifications about changes to a canvas.
// Updates are represented by a bit matrix capturing the new state of the canvas.
type CanvasObserver chan<- *bits.Matrix

// A Canvas represents a simple drawing surface, in which pixels may be either on or off, specifically with dot-matrix displays in mind.
//
// Canvas provides a number of additional benefits over a raw bit-matrix, represented by bits.Matrix.
// * Canvas is thread-safe, and may be read from and written to concurrently.
// * Canvas is observable, sending notifications to observers on change.
// * Allows batch updates, so that notifications are not sent until the batch update is complete.
type Canvas struct {
	buff        *bits.Matrix
	observers   map[uint64]CanvasObserver
	mutex       *sync.RWMutex
	updateDepth uint32
	observerID  uint64
}

// NewCanvas returns a new instance of Canvas, with given dimensions.
func NewCanvas(height, width int) *Canvas {
	return &Canvas{
		buff:      bits.NewMatrix(height, width),
		observers: make(map[uint64]CanvasObserver),
		mutex:     &sync.RWMutex{},
	}
}

// Size returns the size of the canvas.
func (c *Canvas) Size() (height, width int) {
	return c.buff.Size()
}

// Get returns the state of the pixel at the given coodinates.
func (c *Canvas) Get(row, col int) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.buff.Get(row, col)
}

// Set sets the state of the pixel at the given coordinates.
func (c *Canvas) Set(row, col int, value bool) {
	c.mutex.Lock()
	defer func() {
		c.updated()
		c.mutex.Unlock()
	}()
	c.buff.Set(row, col, value)
}

// Clear resets all the pixels in the canvas to off.
func (c *Canvas) Clear() {
	c.mutex.Lock()
	defer func() {
		c.updated()
		c.mutex.Unlock()
	}()
	c.buff.Clear()
}

// Matrix returns a representation of the canvas as a bit-matrix.
// The resulting matrix is by-value; mutating the matrix will not affect the canvas, or vice versa.
func (c *Canvas) Matrix() *bits.Matrix {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.buff.Clone()
}

// Write copies a block of bits from a bit matrix to the canvas at a given location.
// Pixel states are overwritten.  If the source matrix overlaps the extent of the canvas,
// the source matrix is clipped.
// Panics if the destination location is out of bounds.
func (c *Canvas) Write(source *bits.Matrix, row, col int) {
	h, w := source.Size()
	c.mutex.Lock()
	defer func() {
		c.updated()
		c.mutex.Unlock()
	}()
	bits.Copy(source, 0, 0, c.buff, row, col, h, w)
}

// AddObserver registers an observer for this canvas.
// After registration, the observer will receive update notifications when the canvas is modified.
// Returns a unique ID for the observer on this canvas, and a representation of the canvas
// at the point that observation started.
func (c *Canvas) AddObserver(observer CanvasObserver) (id uint64, bits *bits.Matrix) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.observerID++
	id = c.observerID
	c.observers[id] = observer
	bits = c.buff.Clone()
	return
}

// RemoveObserver de-registers an observer from this canvas.
// The argument is the ID provide by the original call to AddObserver.
func (c *Canvas) RemoveObserver(id uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.observers, id)
}

// BeginUpdate notifies that a batch update is starting.
// While a batch update is in progress, observers will not be notified of changes to the canvas.
// Batch updates may be nested, with no update being sent until all nested updates are completed.
func (c *Canvas) BeginUpdate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.updateDepth++
}

// EndUpdate completes a batch update started by BeginUpdate.
// Calling EndUpdate to complete the outermost of a set of nested updates causes observers to be
// notified of the updated canvas.
// EndUpdate panics if called when there is no ongoing batch update.
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
