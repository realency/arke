package max7219

import "periph.io/x/conn/v3/spi"

type Flusher interface {
	Flush()
}

type Flush func()

type Controller interface {
	Flusher
	Shutdown() Flush
	Activate() Flush
	Reset()
	SetDecodeMode(mode byte) Flush
	SetDisplayTest() Flush
	ResetDisplayTest() Flush
	SetIntensity(intensity int) Flush
	SetScanLimit(limit int) Flush
}

type ChipController interface {
	Controller
	SetDigit(digit int, data byte) Flush
}

type ChainController interface {
	Controller
	SetDigit(digit int, data ...byte) Flush
	GetChainLength() int
	SelectChip(index int) ChipController
}

const MaxChainLength = 256 // Arbitrary - but plenty!

func Chain(port spi.Port, chainLen int) (ChainController, error) {
	if chainLen < 0 || chainLen > MaxChainLength {
		panic("Chain length out of bounds")
	}

	var b Bus
	var err error
	if b, err = NewBus(port); err != nil {
		return nil, err
	}

	return NewChain(b, chainLen), nil
}

func Chip(port spi.Port) (ChipController, error) {
	var (
		chain ChainController
		err   error
	)

	if chain, err = Chain(port, 1); err != nil {
		return nil, err
	}

	return chain.SelectChip(0), nil
}
