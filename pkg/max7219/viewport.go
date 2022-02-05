package max7219

import (
	"sync"

	"github.com/realency/arke/pkg/bits"
	"github.com/realency/arke/pkg/display"
	"github.com/realency/arke/pkg/viewport"
)

type ViewPort interface {
	viewport.ViewPort
	SetBrightness(bright byte)
}

type ViewPortUpdateKind byte

type ViewPortUpdate struct {
	kind     ViewPortUpdateKind
	rowDelta int
	colDelta int
}

const (
	DigitZeroAtTop    int = 0 // Digits are indexed from top to bottom, and the least significant bit in a digit register appears at the right
	DigitZeroAtRight  int = 1 // Digits are indexed from right to left, and the least significant bit in a digit register appears at the top
	DigitZeroAtBottom int = 2 // Digits are indexed from bottom to top, and the least significant bit in a digit register appears at the left
	DigitZeroAtLeft   int = 3 // Digits are indexed from left to right, and the least significant bit in a digit register appears at the bottom
)

const (
	BlockZeroAtTop    int = 0 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the top
	BlockZeroAtRight  int = 1 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the right
	BlockZeroAtBottom int = 2 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the bottom
	BlockZeroAtLeft   int = 3 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the left
)

const (
	ViewPortNoOp  ViewPortUpdateKind = 0x00
	ViewPortShift ViewPortUpdateKind = 0x01
)

type viewPort struct {
	canvas                  *display.Canvas
	row, col, height, width int
	bus                     Bus
	chainLength             int
	updates                 chan ViewPortUpdate
	mutex                   *sync.Mutex
}

func newViewPort(bus Bus, chainLength int, blockOrientation, chainOrientation int) ViewPort {
	if blockOrientation != DigitZeroAtBottom {
		panic("Not yet supported")
	}

	if chainOrientation != BlockZeroAtRight {
		panic("Not yet implemented")
	}

	var height, width int
	switch chainOrientation {
	case 0, 2:
		height, width = chainLength*8, 8
	case 1, 3:
		height, width = 8, chainLength*8
	}

	return &viewPort{
		bus:         bus,
		chainLength: chainLength,
		height:      height,
		width:       width,
		updates:     make(chan ViewPortUpdate),
		mutex:       &sync.Mutex{},
	}
}

func (vp *viewPort) broadcast(reg Register, data byte) {
	vp.mutex.Lock()
	for i := 0; i < vp.chainLength; i++ {
		vp.bus.Add(reg, data)
	}
	vp.bus.Send()
	vp.mutex.Unlock()
}

func (vp *viewPort) Attach(canvas *display.Canvas, row, col int) {
	vp.canvas = canvas
	vp.row = row
	vp.col = col

	vp.broadcast(ShutdownRegister, Shutdown)
	vp.broadcast(DisplayTestRegister, NoDisplayTest)
	vp.broadcast(ScanLimitRegister, 0x07)
	vp.broadcast(DecodeModeRegister, DecodeNone)
	for i := 0; i < 8; i++ {
		vp.broadcast(DigitRegister(i), 0x00)
	}
	vp.broadcast(ShutdownRegister, NoShutdown)

	updates := make(chan display.CanvasUpdate, 20)
	go func() {
		for {
			var buff *bits.Matrix
			select {
			case update := <-updates:
				buff = update.Buff
			case update := <-vp.updates:
				vp.row += update.rowDelta
				vp.col += update.colDelta
				buff = vp.canvas.Matrix().Clone()
			}

			for i := 0; i < vp.height; i++ {
				reg := DigitRegister(7 - i)
				r := buff.Reader(vp.row+i, vp.col+vp.width-1, bits.Left, 8*vp.chainLength)
				vp.mutex.Lock()
				for j := 0; j < vp.chainLength; j++ {
					data, e := r.ReadByte()
					if e != nil {
						panic(e)
					}
					vp.bus.Add(reg, data)
				}
				vp.bus.Send()
				vp.mutex.Unlock()
			}
		}
	}()

	canvas.Observe(updates)
}

func (vp *viewPort) SetBrightness(bright byte) {
	vp.broadcast(IntensityRegister, bright)
}

func (vp *viewPort) Shift(rowDelta, colDelta int) {
	vp.updates <- ViewPortUpdate{
		rowDelta: rowDelta,
		colDelta: colDelta,
		kind:     ViewPortShift,
	}
}

func (vp *viewPort) Coords() (row, col int) {
	return vp.row, vp.col
}

func (vp *viewPort) Size() (height, width int) {
	return vp.height, vp.width
}
