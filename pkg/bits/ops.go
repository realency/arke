package bits

// And performs a bitwise and operation.
// The result of the operation is applied to this matrix.
func (m *Matrix) And(other *Matrix) {
	if m.height != other.height || m.width != other.width {
		panic("Mismatched matrix sizes in bits.And")
	}

	for i, v := range m.bits {
		m.bits[i] = v & other.bits[i]
	}
}

// Or performs a bitwise or operation.
// The result of the operation is applied to this matrix.
func (m *Matrix) Or(other *Matrix) {
	if m.height != other.height || m.width != other.width {
		panic("Mismatched matrix sizes in bits.Or")
	}

	for i, v := range m.bits {
		m.bits[i] = v | other.bits[i]
	}
}

// Xor performs a bitwise xor operation.
// The result of the operation is applied to this matrix.
func (m *Matrix) Xor(other *Matrix) {
	if m.height != other.height || m.width != other.width {
		panic("Mismatched matrix sizes in bits.Xor")
	}

	for i, v := range m.bits {
		m.bits[i] = v ^ other.bits[i]
	}
}

// Not performs a bitwise complement operation.
// The result of the operation is applied to this matrix.
func (m *Matrix) Not() {
	for i, v := range m.bits {
		m.bits[i] = v ^ 0xFFFFFFFF
	}
}
