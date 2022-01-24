package max7219

import (
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
)

// Bus provides structured access to a MAX7219 chip, or chain of cascaded chips attached on a serial port.
type Bus interface {
	Write(ops ...Op)
}

// Op represents a single operation on a single MAX7219 chip, setting the state of one register
type Op struct {
	// The address of the register to set
	Register Register
	// The data to set
	Data byte
}

var noOp Op = Op{NoOpRegister, 0x00}

func NoOp() Op {
	return noOp
}

func NewBus(port spi.Port) (Bus, error) {
	var c conn.Conn
	var err error
	if c, err = port.Connect(physic.MegaHertz*10, spi.Mode3, 8); err != nil {
		return nil, err
	}
	return newBus(c), nil
}

type bus struct {
	wire chan []byte
}

func (b *bus) Write(ops ...Op) {
	bytes := make([]byte, len(ops)*2)
	i := 0
	for _, op := range ops {
		bytes[i] = byte(op.Register)
		i++
		bytes[i] = op.Data
		i++
	}
	b.wire <- bytes
}

func newBus(conn conn.Conn) Bus {
	wire := make(chan []byte)

	go func() {
		for {
			conn.Tx(<-wire, nil)
		}
	}()

	return &bus{
		wire: wire,
	}
}
