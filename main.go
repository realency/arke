package main

import "fmt"

func main() {

	b := make([]byte, 10)
	fmt.Println(len(b))
	b = b[0:0]
	fmt.Println(len(b))
}
