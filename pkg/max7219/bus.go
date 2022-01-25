package max7219

import "periph.io/x/conn/v3"

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

type bus struct {
	wire chan []byte
}

func newBus(cx conn.Conn) Bus {
	wire := make(chan []byte)

	go func() {
		for {
			cx.Tx(<-wire, nil)
		}
	}()

	return &bus{
		wire: wire,
	}
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
