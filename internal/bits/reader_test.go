package bits

import (
	"io"
	"testing"
)

func Test_Reader_Reads_To_End_Of_Slice_Where_Target_Is_Longer_And_Offset_Is_Non_Zero(t *testing.T) {
	data := []byte{0x0F, 0x0F, 0x0F}
	r := newReader(data, 3)
	p := make([]byte, 5)
	n, e := r.Read(p)
	if e != nil {
		t.Errorf("Unexpected error")
	}

	if n != 3 {
		t.Errorf("Wrong number of bytes read")
	}

	expected := []byte{0xE1, 0xE1, 0x01, 0x00, 0x00}
	for i := 0; i < 5; i++ {
		if p[i] != expected[i] {
			t.Errorf("Wrong byte value")
		}
	}

	n, e = r.Read(p)
	if n != 0 {
		t.Error("Second read returned non-zero")
	}

	if e != io.EOF {
		t.Error("Second read did not return EOF")
	}
}
