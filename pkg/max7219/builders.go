package max7219

import (
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

// BusBuilder is a builder type for creating a Bus in a fluent programming style.
type BusBuilder struct {
	dev  string
	port spi.Port
	cx   conn.Conn
	bus  Bus
}

// ViewPortBuilder is a builder type for creating a ViewPort in a fluent programming style.
type ViewPortBuilder struct {
	bus              *BusBuilder
	blockOrientation int
	chainOrientation int
	chainLength      int
}

// FromScratch creates a new BusBuilder appropriate for building a bus from scratch using the default SPI device.
func FromScratch() *BusBuilder {
	return FromDeviceName("")
}

// FromDeviceName creates a new BusBuilder appropriate for building a bus using a specific SPI device, identified by device name.
func FromDeviceName(dev string) *BusBuilder {
	return &BusBuilder{
		dev: dev,
	}
}

// FromSpiPort creates a new BusBuilder appropriate for building a bus using a pre-created SPI Port.
func FromSpiPort(port spi.Port) *BusBuilder {
	return &BusBuilder{
		port: port,
	}
}

// FromConnection creates a new BusBuilder appropriate for building a bus using a pre-created Connection.
func FromConnection(cx conn.Conn) *BusBuilder {
	return &BusBuilder{
		cx: cx,
	}
}

// FromBus creates a new BusBuilder that returns a pre-created bus.
func FromBus(bus Bus) *BusBuilder {
	return &BusBuilder{
		bus: bus,
	}
}

// Build builds a Bus using the configuration supplied to the builder, or returns an error if the configuration is incomplete.
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

// WithChainLength specifies the chain length for the bus and returns a ViewPortBuilder.
func (b *BusBuilder) WithChainLength(length int) *ViewPortBuilder {
	return &ViewPortBuilder{
		bus:         b,
		chainLength: length,
	}
}

// WithOrientation specifies the orientation of the blocks in the chains, and of the chain itself and returns the ViewPortBuilder.
func (v *ViewPortBuilder) WithOrientation(blockOrientation, chainOrientation int) *ViewPortBuilder {
	v.blockOrientation = blockOrientation
	v.chainOrientation = chainOrientation
	return v
}

// Build builds the viewport, ready to be attached to a canvas.
func (v *ViewPortBuilder) Build() (*ViewPort, error) {
	b, err := v.bus.Build()
	if err != nil {
		return nil, err
	}

	return newViewPort(b, v.chainLength, v.blockOrientation, v.chainOrientation), nil
}
