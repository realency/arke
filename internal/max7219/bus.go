package max7219

import (
	"realency/arke/pkg/max7219"

	"periph.io/x/conn/v3"
)

type bus struct {
	wire chan []byte
}

func (b *bus) Write(ops ...max7219.Op) {
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

func NewBus(conn conn.Conn) max7219.Bus {
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
