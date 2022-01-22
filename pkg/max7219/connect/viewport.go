package connect

import (
	. "github.com/realency/arke/internal/max7219"
	"github.com/realency/arke/pkg/display"
	"github.com/realency/arke/pkg/max7219"
)

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

func ViewPort(controller max7219.ChainController, blockOrientation, chainOrientation int) display.ViewPort {
	if blockOrientation != DigitZeroAtBottom {
		panic("Not yet supported")
	}

	if chainOrientation != BlockZeroAtRight {
		panic("Not yet implemented")
	}

	return NewViewPort(controller, blockOrientation, chainOrientation)
}
