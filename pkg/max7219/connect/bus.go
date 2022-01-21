package connect

import (
	internal "realency/arke/internal/max7219"
	"realency/arke/pkg/max7219"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
)

func Bus(port spi.Port) (max7219.Bus, error) {
	var c conn.Conn
	var err error
	if c, err = port.Connect(physic.MegaHertz*10, spi.Mode3, 8); err != nil {
		return nil, err
	}
	return internal.NewBus(c), nil
}
