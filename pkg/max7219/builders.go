package max7219

import (
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

type BusBuilder struct {
	dev  string
	port spi.Port
	cx   conn.Conn
	bus  Bus
}

type ChainBuilder struct {
	bus    *BusBuilder
	length int
	chain  ChainController
}

type ViewPortBuilder struct {
	chain            *ChainBuilder
	blockOrientation int
	chainOrientation int
}

func FromScratch() *BusBuilder {
	return FromDeviceName("")
}

func FromDeviceName(dev string) *BusBuilder {
	return &BusBuilder{
		dev: dev,
	}
}

func FromSpiPort(port spi.Port) *BusBuilder {
	return &BusBuilder{
		port: port,
	}
}

func FromConnection(cx conn.Conn) *BusBuilder {
	return &BusBuilder{
		cx: cx,
	}
}

func FromBus(bus Bus) *BusBuilder {
	return &BusBuilder{
		bus: bus,
	}
}

func FromChain(chain ChainController) *ChainBuilder {
	return &ChainBuilder{
		chain: chain,
	}
}

func (b *BusBuilder) Build() (Bus, error) {
	var err error

	if b.bus != nil {
		return b.bus, nil
	}

	if b.cx != nil {
		return newBus(b.cx), nil
	}

	if b.port == nil {
		host.Init()
		if b.port, err = spireg.Open(b.dev); err != nil {
			return nil, err
		}
	}

	if b.cx, err = b.port.Connect(physic.MegaHertz*10, spi.Mode3, 8); err != nil {
		return nil, err
	}

	return newBus(b.cx), nil
}

func (b *BusBuilder) WithChainLength(length int) *ChainBuilder {
	return &ChainBuilder{
		bus:    b,
		length: length,
	}
}

func (c *ChainBuilder) Build() (ChainController, error) {
	if c.chain != nil {
		return c.chain, nil
	}

	b, err := c.bus.Build()
	if err != nil {
		return nil, err
	}

	return newChain(b, c.length), nil
}

func (c *ChainBuilder) WithOrientation(blockOrientation, chainOrientation int) *ViewPortBuilder {
	return &ViewPortBuilder{
		chain:            c,
		blockOrientation: blockOrientation,
		chainOrientation: chainOrientation,
	}
}

func (v *ViewPortBuilder) Build() (ViewPort, error) {
	c, err := v.chain.Build()
	if err != nil {
		return nil, err
	}

	return newViewPort(c, v.blockOrientation, v.chainOrientation), nil
}
