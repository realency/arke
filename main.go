package main

import (
	"fmt"

	"github.com/realency/arke/internal/bits"
)

func main() {
	b := bits.NewBuffer(32, 32)

	b.Set(10, 10, true)
	b.Set(10, 30, true)
	b.Set(10, 31, true)

	for i := 0; i < 32; i++ {
		s := ""
		for j := 0; j < 32; j++ {
			if b.Get(i, j) {
				s += "@ "
			} else {
				s += ". "
			}
		}
		fmt.Println(s)
	}

	d := make([]byte, 1)
	r := b.RowReader(10, 18)
	count, _ := r.Read(d)
	fmt.Println(count)
	s := ""
	for i := 0; i < count; i++ {
		s += fmt.Sprintf("%d ", d[i])
	}
	fmt.Println(s)
}
