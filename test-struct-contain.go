package main

import "fmt"

type a struct {
}

func (a *a) IsStruct() bool {
	return true
}

type b struct {
	a
}

func (b *b) IsStruct() bool {
	return false
}

func main() {
	a := new(a)
	b := new(b)

	fmt.Println(a.IsStruct())
	fmt.Println(b.IsStruct())
}
