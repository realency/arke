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

type ViewPortBuilder struct {
	bus              *BusBuilder
	blockOrientation int
	chainOrientation int
	chainLength      int
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

func (b *BusBuilder) WithChainLength(length int) *ViewPortBuilder {
	return &ViewPortBuilder{
		bus:         b,
		chainLength: length,
	}
}

func (v *ViewPortBuilder) WithOrientation(blockOrientation, chainOrientation int) *ViewPortBuilder {
	v.blockOrientation = blockOrientation
	v.chainOrientation = chainOrientation
	return v
}

func (v *ViewPortBuilder) Build() (*ViewPort, error) {
	b, err := v.bus.Build()
	if err != nil {
		return nil, err
	}

	return newViewPort(b, v.chainLength, v.blockOrientation, v.chainOrientation), nil
}
