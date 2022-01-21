package max7219

type Flusher interface {
	Flush()
}

type Flush func()

type Controller interface {
	Flusher
	Shutdown() Flush
	Activate() Flush
	Reset()
	SetDecodeMode(mode byte) Flush
	SetDisplayTest() Flush
	ResetDisplayTest() Flush
	SetIntensity(intensity int) Flush
	SetScanLimit(limit int) Flush
}

type ChipController interface {
	Controller
	SetDigit(digit int, data byte) Flush
}

type ChainController interface {
	Controller
	SetDigit(digit int, data ...byte) Flush
	GetChainLength() int
	SelectChip(index int) ChipController
}

const MaxChainLength = 256 // Arbitrary - but plenty!
