package max7219

type chain struct {
	bus      Bus
	chainLen int
	chips    []*chip
	buff     []Op
}

func NewChain(bus Bus, chainLen int) ChainController {
	result := &chain{
		bus:      bus,
		buff:     make([]Op, chainLen),
		chainLen: chainLen,
	}

	result.Reset()
	return result
}

func (c *chain) Activate() Flush {
	c.set(ShutdownRegister, NoShutdown)
	return c.getFlush(ShutdownRegister)
}

func (c *chain) Flush() {
	c.flush(DisplayTestRegister)
	for r := Register(0x01); r <= Register(0x0C); r++ {
		c.flush(r)
	}
}

func (c *chain) GetChainLength() int {
	return len(c.chips)
}

func (c *chain) Reset() {
	c.send(ShutdownRegister, Shutdown)
	c.send(DisplayTestRegister, NoDisplayTest)
	for b := byte(0x01); b <= 0x0A; b++ {
		c.send(Register(b), 0x00)
	}
	c.send(ScanLimitRegister, 0x07)
	c.send(ShutdownRegister, NoShutdown)

	c.chips = make([]*chip, c.chainLen)
	for i := range c.chips {
		c.chips[i] = newChip(c, i)
	}
}

func (c *chain) SetDigit(digit int, data ...byte) Flush {
	reg := DigitRegister(digit)
	switch len(data) {
	case 1:
		for _, chip := range c.chips {
			chip.set(reg, data[0])
		}
	case c.chainLen:
		for i, chip := range c.chips {
			chip.set(reg, data[i])
		}
	default:
		panic("Must send exactly one byte for each chip, or a single byte to be repeated across all chips")
	}
	return c.getFlush(reg)
}

func (c *chain) SetDecodeMode(mode byte) Flush {
	c.set(DecodeModeRegister, mode)
	return c.getFlush(DecodeModeRegister)
}

func (c *chain) SetIntensity(intensity int) Flush {
	c.set(IntensityRegister, Intensity(intensity))
	return c.getFlush(IntensityRegister)
}

func (c *chain) Shutdown() Flush {
	c.set(ShutdownRegister, Shutdown)
	return c.getFlush(ShutdownRegister)
}

func (c *chain) SetDisplayTest() Flush {
	c.set(DisplayTestRegister, DisplayTest)
	return c.getFlush(DisplayTestRegister)
}

func (c *chain) ResetDisplayTest() Flush {
	c.set(DisplayTestRegister, NoDisplayTest)
	return c.getFlush(DisplayTestRegister)
}

func (c *chain) SetScanLimit(limit int) Flush {
	c.set(ScanLimitRegister, ScanLimit(limit))
	return c.getFlush(ScanLimitRegister)
}

func (c *chain) SelectChip(index int) ChipController {
	if index < 0 || index >= len(c.chips) {
		panic("Chip index out of range")
	}
	return c.chips[index]
}

func (c *chain) set(reg Register, data byte) {
	for _, chip := range c.chips {
		chip.set(reg, data)
	}
}

func (c *chain) getFlush(reg Register) Flush {
	return func() {
		c.flush(reg)
	}
}

func (c *chain) send(reg Register, data byte) {
	op := Op{
		Register: reg,
		Data:     data,
	}

	for i := range c.buff {
		c.buff[i] = op
	}

	c.bus.Write(c.buff...)
}

func (c *chain) flush(reg Register) {
	mask := uint16(0x0001) << reg
	dirty := false
	for i, chip := range c.chips {
		if chip.dirtyFlags&mask == 0 {
			c.buff[i] = Op{}
			continue
		}

		dirty = true
		c.buff[i] = Op{
			Register: reg,
			Data:     chip.registers[reg],
		}
		chip.dirtyFlags &= ^mask
	}

	if dirty {
		c.bus.Write(c.buff...)
	}
}
