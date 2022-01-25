package max7219

import (
	"github.com/realency/arke/pkg/viewport"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

type BusBuilder interface {
	Build() (Bus, error)
	WithChainLength(length int) ChainBuilder
}

type ChainBuilder interface {
	Build() (ChainController, error)
	WithOrientation(blockOrientation, chainOrientation int) ViewPortBuilder
}

type ViewPortBuilder interface {
	Build() (viewport.ViewPort, error)
}

type busBuilder struct {
	dev  string
	port spi.Port
	cx   conn.Conn
	bus  Bus
}

type chainBuilder struct {
	bus    *busBuilder
	length int
	chain  ChainController
}

type viewportBuilder struct {
	chain            *chainBuilder
	blockOrientation int
	chainOrientation int
}

func FromScratch() BusBuilder {
	return FromDeviceName("")
}

func FromDeviceName(dev string) BusBuilder {
	return &busBuilder{
		dev: dev,
	}
}

func FromSpiPort(port spi.Port) BusBuilder {
	return &busBuilder{
		port: port,
	}
}

func FromConnection(cx conn.Conn) BusBuilder {
	return &busBuilder{
		cx: cx,
	}
}

func FromBus(bus Bus) BusBuilder {
	return &busBuilder{
		bus: bus,
	}
}

func FromChain(chain ChainController) ChainBuilder {
	return &chainBuilder{
		chain: chain,
	}
}

func (b *busBuilder) Build() (Bus, error) {
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

func (b *busBuilder) WithChainLength(length int) ChainBuilder {
	return &chainBuilder{
		bus:    b,
		length: length,
	}
}

func (c *chainBuilder) Build() (ChainController, error) {
	if c.chain != nil {
		return c.chain, nil
	}

	b, err := c.bus.Build()
	if err != nil {
		return nil, err
	}

	return newChain(b, c.length), nil
}

func (c *chainBuilder) WithOrientation(blockOrientation, chainOrientation int) ViewPortBuilder {
	return &viewportBuilder{
		chain:            c,
		blockOrientation: blockOrientation,
		chainOrientation: chainOrientation,
	}
}

func (v *viewportBuilder) Build() (viewport.ViewPort, error) {
	c, err := v.chain.Build()
	if err != nil {
		return nil, err
	}

	return newViewPort(c, v.blockOrientation, v.chainOrientation), nil
}
