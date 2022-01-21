package max7219

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
