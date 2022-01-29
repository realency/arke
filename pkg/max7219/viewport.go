package max7219

import (
	"github.com/realency/arke/internal/bits"
	"github.com/realency/arke/pkg/display"
	"github.com/realency/arke/pkg/viewport"
)

type ViewPort interface {
	viewport.ViewPort
	Chain() ChainController
}

const (
	DigitZeroAtTop    int = 0
	DigitZeroAtRight  int = 1
	DigitZeroAtBottom int = 2
	DigitZeroAtLeft   int = 3
)

const (
	BlockZeroAtTop    int = 0
	BlockZeroAtRight  int = 1
	BlockZeroAtBottom int = 2
	BlockZeroAtLeft   int = 3
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
	updates := make(chan *bits.ImmutableBuffer)
	go func() {
		for {
			bits := <-updates
			buff := make([]byte, vp.controller.GetChainLength())
			reversed := make([]byte, vp.controller.GetChainLength())
			for i := 0; i < vp.height; i++ {
				bits.RowReader(i+row, col).Read(buff)
				for j, b := range buff {
					reversed[vp.controller.GetChainLength()-(j+1)] = b
				}
				vp.controller.SetDigit(7-i, reversed...)
			}
			vp.controller.Flush()
		}
	}()

	canvas.Observe(updates)
}

func (vp *viewPort) Chain() ChainController {
	return vp.controller
}
