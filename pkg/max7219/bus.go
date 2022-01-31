package max7219

import "periph.io/x/conn/v3"

// Bus provides structured access to a MAX7219 chip, or chain of cascaded chips attached on a serial port.
type Bus interface {
	Add(reg Register, data byte)
	Send()
}

type bus struct {
	buff []byte
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
		buff: make([]byte, 0, 1024),
		wire: wire,
	}
}

func (b *bus) Add(reg Register, data byte) {
	b.buff = append(b.buff, byte(reg), data)
}

func (b *bus) Send() {
	b.wire <- b.buff
	b.buff = b.buff[0:0]
}
