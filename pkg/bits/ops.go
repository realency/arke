package bits

func (m *Matrix) And(other *Matrix) {
	if m.height != other.height || m.width != other.width {
		panic("Mismatched matrix sizes in bits.And")
	}

	for i, v := range m.bits {
		m.bits[i] = v & other.bits[i]
	}
}

func (m *Matrix) Or(other *Matrix) {
	if m.height != other.height || m.width != other.width {
		panic("Mismatched matrix sizes in bits.Or")
	}

	for i, v := range m.bits {
		m.bits[i] = v | other.bits[i]
	}
}

func (m *Matrix) Xor(other *Matrix) {
	if m.height != other.height || m.width != other.width {
		panic("Mismatched matrix sizes in bits.Xor")
	}

	for i, v := range m.bits {
		m.bits[i] = v ^ other.bits[i]
	}
}
