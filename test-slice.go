package main

import (
	"fmt"
)

func main() {
	x := make([]byte, 10)

	fmt.Println(len(x), cap(x))

	y := x[:3]
	fmt.Println(len(y), cap(y))
	fmt.Println(len(x), cap(x))

	z := x[10:]
	fmt.Println(len(z), cap(y))
	fmt.Println(len(x), cap(x))

	c := make(chan bool)

	if c == nil {
		fmt.Println("c channel is nil")
	} else {
		fmt.Println("c channel is not nil")
	}

	close(c)

	if c == nil {
		fmt.Println("c channel is nil")
	} else {
		fmt.Println("c channel is not nil")
	}
}
