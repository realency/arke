package max7219

// Register represents the address of one of the registers of the MAX7219 chip.
type Register byte

// Constant definitions of the register addresses.
const (
	NoOpRegister        Register = 0x00
	Digit0Register      Register = 0x01
	Digit1Register      Register = 0x02
	Digit2Register      Register = 0x03
	Digit3Register      Register = 0x04
	Digit4Register      Register = 0x05
	Digit5Register      Register = 0x06
	Digit6Register      Register = 0x07
	Digit7Register      Register = 0x08
	DecodeModeRegister  Register = 0x09
	IntensityRegister   Register = 0x0A
	ScanLimitRegister   Register = 0x0B
	ShutdownRegister    Register = 0x0C
	DisplayTestRegister Register = 0x0F
)

// Constant definitions for decode mode values.
const (
	DecodeNone byte = 0x00
	DecodeAll  byte = 0xFF
)

// Constant definitions for BCD Code B character codes, used when Decode mode is on.
const (
	Char0     byte = 0x00
	Char1     byte = 0x01
	Char2     byte = 0x02
	Char3     byte = 0x03
	Char4     byte = 0x04
	Char5     byte = 0x05
	Char6     byte = 0x06
	Char7     byte = 0x07
	Char8     byte = 0x08
	Char9     byte = 0x09
	CharDash  byte = 0x0A
	CharE     byte = 0x0B
	CharH     byte = 0x0C
	CharL     byte = 0x0D
	CharP     byte = 0x0E
	CharBlank byte = 0x0F
)

// Constant definitions for Shutdown mode data values.
const (
	Shutdown   byte = 0x00
	NoShutdown byte = 0x01
)

// Constant definitions for Display Test mode data values.
const (
	DisplayTest   byte = 0x01
	NoDisplayTest byte = 0x00
)

// DigitRegister returns the register address for a given digit's index.
// The digit register controls either a single digit in a 7-segment display,
// or a line of LEDs in a dot-matrix display.
//
// Panics if the argument is out of range.  Acceptable argument values are 0..7, inclusive.
func DigitRegister(digit int) Register {
	if digit < 0 || digit > 7 {
		panic("digit register out of range. The MAX7219 chip has 8 digit registers (in the range 0..7)")
	}
	return Register(digit + 1)
}

// Intensity returns the data value for the intensity of the display, given a numerical value between 0 and 15, inclusively.
//
// Panics if the argument is out of range.
func Intensity(intensity int) byte {
	if intensity < 0 || intensity > 15 {
		panic("Intensity value out of range. Intensity must be in the range 0..15, inclusively")
	}

	return byte(intensity)
}

// ScanLimit sets a limit to the number of digits displayed.
//
// Using a limit argument of 0 causes only the first digit (index 0) to be displayed.  Using a limit argument of 1 causes the first and second digits to be displayed. And so forth.
// All digits will be displayed for an argument of 7.  For a dot-matrix display, digits map to lines in the display, so scan limit affects the number of lines displayed.
func ScanLimit(limit int) byte {
	if limit < 0 || limit > 15 {
		panic("Scan limit must be in the range 0..7, inclusively")
	}
	return byte(limit)
}

// CodedChar returns the BCD Code B character code for a given rune.
//
// For use on a 7-segment display when in Decode mode.
func CodedChar(from rune, decimalPoint bool) byte {
	var result byte
	if decimalPoint {
		result = 0x08
	}

	switch {
	case from >= '0' && from <= '9':
		return result | byte(from-'0')
	case from == 'h' || from == 'H':
		return result | CharH
	case from == 'e' || from == 'E':
		return result | CharE
	case from == 'l' || from == 'L':
		return result | CharL
	case from == 'p' || from == 'P':
		return result | CharP
	case from == '-':
		return result | CharDash
	case from == ' ' || from == '\u0000':
		return result | CharBlank
	}

	panic("Character is not representable in BCD Code B")
}
