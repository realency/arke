package max7219

import (
	"github.com/realency/arke/pkg/bits"
	"github.com/realency/arke/pkg/display"
	"github.com/realency/arke/pkg/viewport"
)

type ViewPort interface {
	viewport.ViewPort
	Chain() ChainController
}

const (
	DigitZeroAtTop    int = 0 // Digits are indexed from top to bottom, and the least significant bit in a digit register appears at the left
	DigitZeroAtRight  int = 1 // Digits are indexed from right to left, and the least significant bit in a digit register appears at the top
	DigitZeroAtBottom int = 2 // Digits are indexed from bottom to top, and the least significant bit in a digit register appears at the right
	DigitZeroAtLeft   int = 3 // Digits are indexed from left to right, and the least significant bit in a digit register appears at the bottom
)

const (
	BlockZeroAtTop    int = 0 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the top
	BlockZeroAtRight  int = 1 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the right
	BlockZeroAtBottom int = 2 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the bottom
	BlockZeroAtLeft   int = 3 // In a chain of chips, the block controlled by the first address-byte pair sent in a packet is at the left
)

type viewPort struct {
	controller    ChainController
	height, width int
}

func newViewPort(controller ChainController, blockOrientation, chainOrientation int) ViewPort {
	if blockOrientation != DigitZeroAtBottom {
		panic("Not yet supported")
	}

	if chainOrientation != BlockZeroAtRight {
		panic("Not yet implemented")
	}

	var height, width int
	switch chainOrientation {
	case 0, 2:
		height, width = controller.GetChainLength()*8, 8
	case 1, 3:
		height, width = 8, controller.GetChainLength()*8
	}

	return &viewPort{
		controller: controller,
		height:     height,
		width:      width,
	}
}

func (vp *viewPort) Attach(canvas *display.Canvas, row, col int) {
	updates := make(chan display.CanvasUpdate, 20)
	go func() {
		for {
			update := <-updates
			buff := make([]byte, vp.controller.GetChainLength())
			for i := 0; i < vp.height; i++ {
				r := update.Buff.Reader(row+i, col, bits.Right, 8*vp.controller.GetChainLength())
				for j := 0; j < len(buff); j++ {
					b, e := r.ReadByte()
					if e != nil {
						panic(e)
					}
					buff[j] = b
				}
				vp.controller.SetDigit(7-i, buff...)
			}
			vp.controller.Flush()
		}
	}()

	canvas.Observe(updates)
}

func (vp *viewPort) Chain() ChainController {
	return vp.controller
}
