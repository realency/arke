package max7219

import (
	"github.com/realency/arke/pkg/bits"
	"github.com/realency/arke/pkg/display"
	"github.com/realency/arke/pkg/viewport"
)

type ViewPortUpdateKind byte

type offset struct {
	rowDelta int
	colDelta int
}

const (
	DigitZeroAtTop    int = 0 // Digits are indexed from top to bottom, and the least significant bit in a digit register controls an LED at the right
	DigitZeroAtRight  int = 1 // Digits are indexed from right to left, and the least significant bit in a digit register controls an LED at the bottom
	DigitZeroAtBottom int = 2 // Digits are indexed from bottom to top, and the least significant bit in a digit register controls an LED at the left
	DigitZeroAtLeft   int = 3 // Digits are indexed from left to right, and the least significant bit in a digit register controls an LED at the top
)

const (
	BlockZeroAtTop    int = 0 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the top
	BlockZeroAtRight  int = 1 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the right
	BlockZeroAtBottom int = 2 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the bottom
	BlockZeroAtLeft   int = 3 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the left
)

type ViewPort struct {
	canvas                  *display.Canvas
	row, col, height, width int
	bus                     Bus
	chainLength             int
	brightness              chan byte
	shifts                  chan offset
}

func newViewPort(bus Bus, chainLength int, blockOrientation, chainOrientation int) *ViewPort {
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

	result := &ViewPort{
		bus:         bus,
		chainLength: chainLength,
		height:      height,
		width:       width,
	}

	result.init()
	return result
}

func (vp *ViewPort) init() {
	vp.broadcast(ShutdownRegister, Shutdown)
	vp.broadcast(DisplayTestRegister, NoDisplayTest)
	vp.broadcast(ScanLimitRegister, 0x07)
	vp.broadcast(DecodeModeRegister, DecodeNone)
	for i := 0; i < 8; i++ {
		vp.broadcast(DigitRegister(i), 0x00)
	}
	vp.broadcast(ShutdownRegister, NoShutdown)

	canvasUpdates := make(chan display.CanvasUpdate, 20)

	go func() {
		for {
			select {
			case update := <-canvasUpdates:
				vp.handleUpdate(update.Buff)
			case b := <-vp.brightness:
				vp.broadcast(IntensityRegister, b)
			case shift := <-vp.shifts:
				vp.row += shift.rowDelta
				vp.col += shift.colDelta
				vp.handleUpdate(vp.canvas.Matrix().Clone())
			}
		}

	}()
}

func (vp *ViewPort) broadcast(reg Register, data byte) {
	for i := 0; i < vp.chainLength; i++ {
		vp.bus.Add(reg, data)
	}
	vp.bus.Send()
}

func (vp *ViewPort) Attach(canvas *display.Canvas, row, col int) error {
	if vp.canvas != nil {
		return viewport.EAlreadyAttached
	}

	vp.canvas = canvas
	vp.row = row
	vp.col = col

	vp.id = canvas.AddObserver(canvasUpdates)
	return nil
}

func (vp *ViewPort) Detach() error {
	if vp.canvas == nil {
		return viewport.ENotAttached
	}

	vp.canvas.RemoveObserver(vp.id)

	return nil
}

func (vp *ViewPort) handleUpdate(buff *bits.Matrix) {
	for i := 0; i < vp.height; i++ {
		reg := DigitRegister(7 - i)
		r := buff.Reader(vp.row+i, vp.col+vp.width-1, bits.Left, 8*vp.chainLength)
		for j := 0; j < vp.chainLength; j++ {
			data, e := r.ReadByte()
			if e != nil {
				panic(e)
			}
			vp.bus.Add(reg, data)
		}
		vp.bus.Send()
	}
}

func (vp *ViewPort) SetBrightness(bright byte) {
	vp.updates <- ViewPortUpdate{
		kind:       ViewPortBrightness,
		brightness: bright,
	}
}

func (vp *ViewPort) Shift(rowDelta, colDelta int) {
	vp.updates <- ViewPortUpdate{
		rowDelta: rowDelta,
		colDelta: colDelta,
		kind:     ViewPortShift,
	}
}

func (vp *ViewPort) Offset() (row, col int) {
	return vp.row, vp.col
}

func (vp *ViewPort) Size() (height, width int) {
	return vp.height, vp.width
}
