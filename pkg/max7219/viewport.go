package max7219

import (
	"log"

	"github.com/realency/arke/pkg/bits"
	"github.com/realency/arke/pkg/display"
)

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

type offset struct {
	row, col int
}

type attachment struct {
	offset offset
	canvas *display.Canvas
}

type ViewPort struct {
	canvas                  *display.Canvas
	id                      uint64
	row, col, height, width int
	bus                     Bus
	chainLength             int
	offsets                 chan offset
	brightness              chan byte
	canvasUpdates           chan *bits.Matrix
	attachments             chan attachment
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
		bus:           bus,
		chainLength:   chainLength,
		height:        height,
		width:         width,
		offsets:       make(chan offset),
		brightness:    make(chan byte, 20),
		canvasUpdates: make(chan *bits.Matrix, 20),
		attachments:   make(chan attachment, 20),
	}

	result.init()
	go result.run()
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

}

func (vp *ViewPort) run() {
	for {
		if len(vp.canvasUpdates) > 10 || len(vp.brightness) > 10 || len(vp.attachments) > 10 || len(vp.offsets) > 10 {
			log.Println("WARNING ViewPort buffering operations")
		}

		if len(vp.canvasUpdates) == 20 || len(vp.brightness) == 20 || len(vp.attachments) == 20 || len(vp.offsets) == 20 {
			panic("ViewPort buffer overflow")
		}

		select {
		case c := <-vp.canvasUpdates:
			vp.handleUpdate(c)
		case b := <-vp.brightness:
			vp.broadcast(IntensityRegister, b)
		case o := <-vp.offsets:
			if vp.canvas == nil {
				continue
			}
			vp.setOffset(o)
			vp.handleUpdate(vp.canvas.Matrix().Clone())
		case a := <-vp.attachments:
			if vp.canvas == a.canvas {
				continue
			}
			if vp.canvas != nil {
				vp.canvas.RemoveObserver(vp.id)
				vp.id = 0
				vp.row = -1
				vp.col = -1
				vp.canvas = nil
			}
			if a.canvas != nil {
				var b *bits.Matrix
				vp.canvas = a.canvas
				vp.id, b = a.canvas.AddObserver(vp.canvasUpdates)
				vp.setOffset(a.offset)
				vp.handleUpdate(b)
			}
		}
	}
}

func (vp *ViewPort) setOffset(o offset) {
	h, w := vp.canvas.Size()
	if o.row < 0 {
		o.row = 0
	}
	if o.row+vp.height > h {
		o.row = h - vp.height
	}
	if o.col < 0 {
		o.col = 0
	}
	if o.col+vp.width > w {
		o.col = w - vp.width
	}
	vp.row = o.row
	vp.col = o.col
}

func (vp *ViewPort) broadcast(reg Register, data byte) {
	for i := 0; i < vp.chainLength; i++ {
		vp.bus.Add(reg, data)
	}
	vp.bus.Send()
}

func (vp *ViewPort) handleUpdate(buff *bits.Matrix) {
	for i := 0; i < vp.height; i++ {
		reg := DigitRegister(7 - i)

		c := bits.NewCursor(buff, vp.row+i, vp.col+vp.width)
		for j := 0; j < vp.chainLength; j++ {
			data, _ := c.ReadLeftByte()
			vp.bus.Add(reg, data)
		}
		vp.bus.Send()
	}
}

func (vp *ViewPort) Attach(canvas *display.Canvas, row, col int) {
	vp.attachments <- attachment{
		canvas: canvas,
		offset: offset{row, col},
	}
}

func (vp *ViewPort) Detach() {
	vp.attachments <- attachment{canvas: nil}
}

func (vp *ViewPort) SetBrightness(bright byte) {
	if bright > 15 {
		bright = 15
	}
	vp.brightness <- bright
}

func (vp *ViewPort) Locate(row, col int) {
	vp.offsets <- offset{row, col}
}

func (vp *ViewPort) Offset() (row, col int) {
	return vp.row, vp.col
}

func (vp *ViewPort) Size() (height, width int) {
	return vp.height, vp.width
}
