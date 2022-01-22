package max7219

import . "github.com/realency/arke/pkg/max7219"

type chip struct {
	bus        Bus
	chainLen   int
	index      int
	registers  []byte
	dirtyFlags uint16
}

func newChip(chain *chain, index int) *chip {
	result := &chip{
		index:      index,
		chainLen:   len(chain.chips),
		bus:        chain.bus,
		registers:  make([]byte, 16),
		dirtyFlags: 0x0000,
	}
	result.registers[ScanLimitRegister] = 0x07
	return result
}

func (c *chip) SetDigit(digit int, data byte) Flush {
	r := DigitRegister(digit)
	c.set(r, data)
	return c.getFlush(r)
}

func (c *chip) Shutdown() Flush {
	c.set(ShutdownRegister, Shutdown)
	return c.getFlush(ShutdownRegister)
}

func (c *chip) Activate() Flush {
	c.set(ShutdownRegister, NoShutdown)
	return c.getFlush(ShutdownRegister)
}

func (c *chip) SetScanLimit(limit int) Flush {
	c.set(ScanLimitRegister, ScanLimit(limit))
	return c.getFlush(ScanLimitRegister)
}

func (c *chip) SetIntensity(intensity int) Flush {
	c.set(IntensityRegister, Intensity(intensity))
	return c.getFlush(IntensityRegister)
}

func (c *chip) SetDisplayTest() Flush {
	c.set(DisplayTestRegister, DisplayTest)
	return c.getFlush(DisplayTestRegister)
}

func (c *chip) ResetDisplayTest() Flush {
	c.set(DisplayTestRegister, NoDisplayTest)
	return c.getFlush(DisplayTestRegister)
}

func (c *chip) Reset() {
	c.send(ShutdownRegister, Shutdown)
	c.send(DisplayTestRegister, NoDisplayTest)
	for b := byte(0x01); b <= 0x0A; b++ {
		c.send(Register(b), 0x00)
	}
	c.send(ScanLimitRegister, 0x07)
	c.registers = make([]byte, 16)
	c.registers[ScanLimitRegister] = 0x07
	c.send(ShutdownRegister, NoShutdown)
	c.dirtyFlags = 0x0000
}

func (c *chip) SetDecodeMode(mode byte) Flush {
	c.set(DecodeModeRegister, mode)
	return c.getFlush(DecodeModeRegister)
}

func (c *chip) Flush() {
	mask := uint16(0x0002)
	for r := Register(0x01); r <= Register(0x0F); r++ {
		if c.dirtyFlags&mask != 0 {
			c.send(r, c.registers[r])
		}
		mask <<= 1
	}
	c.dirtyFlags = 0x0000
}

func (c *chip) send(reg Register, data byte) {
	packet := make([]Op, c.chainLen)
	packet[c.index] = Op{
		Register: reg,
		Data:     data,
	}
	c.bus.Write(packet...)
}

func (c *chip) flush(reg Register) {
	mask := uint16(0x0001) << reg
	if c.dirtyFlags&mask == 0 {
		return
	}
	c.send(reg, c.registers[reg])
	c.dirtyFlags &= ^mask
}

func (c *chip) set(register Register, data byte) {
	dirtyFlagMask := uint16(0x0001) << register
	if c.registers[register] != data {
		c.registers[register] = data
		c.dirtyFlags |= dirtyFlagMask
	}
}

func (c *chip) getFlush(register Register) Flush {
	return func() {
		c.flush(register)
	}
}
