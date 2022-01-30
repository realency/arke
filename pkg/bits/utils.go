package bits

func Reverse(b byte) byte {
	var result byte = 0x00
	for i := 0; i < 8; i++ {
		result <<= 1
		if (b & 0x01) != 0 {
			result |= 0x01
		}
		b >>= 1
	}
	return result
}
