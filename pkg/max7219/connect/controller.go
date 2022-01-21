package connect

import (
	internal "realency/arke/internal/max7219"
	"realency/arke/pkg/max7219"

	"periph.io/x/conn/v3/spi"
)

func Chain(port spi.Port, chainLen int) (max7219.ChainController, error) {
	if chainLen < 0 || chainLen > max7219.MaxChainLength {
		panic("Chain length out of bounds")
	}

	var b max7219.Bus
	var err error
	if b, err = Bus(port); err != nil {
		return nil, err
	}

	return internal.NewChain(b, chainLen), nil
}

func Chip(port spi.Port) (max7219.ChipController, error) {
	var (
		chain max7219.ChainController
		err   error
	)

	if chain, err = Chain(port, 1); err != nil {
		return nil, err
	}

	return chain.SelectChip(0), nil
}
