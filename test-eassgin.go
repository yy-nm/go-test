package main

import (
	_ "bufio"
	"fmt"
	"os"
)

func init() {

	fmt.Println("hello world!")
}

type none struct {
}

func (n *none) init() {
	fmt.Println("none hello world!")
}

func main() {
	var n *none
	if n != nil {

	}
	f, err := os.Open("test.iml")
	if err != nil {
		fmt.Println("open file err: ", err.Error())
	}

	s, err := f.Stat()
	if err != nil {
		fmt.Println("stat file err: ", err.Error())
	}

	fmt.Println(s)
	//err := f.Close()
	f.Close()
}
