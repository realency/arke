package max7219

import (
	"github.com/realency/arke/pkg/display"

	. "github.com/realency/arke/pkg/max7219"
)

type viewPort struct {
	controller    ChainController
	height, width int
}

func NewViewPort(controller ChainController, blockOrientation, chainOrientation int) display.ViewPort {
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
	observer := func(bits [][]bool) {
		buff := make([]byte, vp.controller.GetChainLength())

		for i := 0; i < vp.height; i++ {
			vp.rowToBytes(bits[i], col, buff)
			vp.controller.SetDigit(7-i, buff...)
		}
		vp.controller.Flush()
	}

	canvas.Observe(observer)
}

func (vp *viewPort) rowToBytes(row []bool, offset int, buff []byte) {
	b := byte(0x00)
	i := vp.controller.GetChainLength() - 1
	m := byte(0x01)

	for {
		if row[offset] {
			b |= m
		}

		m <<= 1

		if m == 0 {
			buff[i] = b
			b = 0x00
			m = 0x01
			i--
			if i < 0 {
				break
			}
		}

		offset++
		if offset >= len(row) {
			break
		}
	}
}
