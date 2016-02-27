package main

import (
	"fmt"
	"test/sf/net"
	"time"
)

func main() {
	x := make([]byte, 10)
	fmt.Println(x[0:2])
	fmt.Println(x[2:3])
	fmt.Println(x[2:])

	y := make([]byte, 0, 10)
	if y == nil {
		fmt.Println("y is nil")
	}
	fmt.Println(y)
	y = append(y, byte(1))
	fmt.Println(y)
	y = append(y, byte(1))
	fmt.Println(y)
	y = append(y, byte(1))
	fmt.Println(y)
	y = append(y, byte(1))
	var z []byte
	copy(y, z)
	fmt.Println(y)
	fmt.Println(len(y))
	fmt.Println(cap(y))
	fmt.Println(y[:])

	z = x[:2]
	fmt.Println(len(z))
	fmt.Println(cap(z))
	z = z[:cap(z)]
	fmt.Println(len(z))
	fmt.Println(cap(z))

	ch := make(chan bool, 1)
	//	ch := make(chan bool)

	//	ch <- true
	go func() {
		time.Sleep(time.Second * 4)
		fmt.Println("sleep done")
		//		v, err := <- ch
		for v := range ch {

			fmt.Println(v)
		}
		//		ch <- true
	}()

	time.Sleep(time.Second * 2)
	fmt.Println("done")
	ch <- true
	fmt.Println("done!")
	time.Sleep(time.Second * 1)
	close(ch)
	//	ch <- true

	time.Sleep(time.Second * 2)
	net.NewPacket(net.MsgType(0), nil, nil)

	x = []byte{1, 2, 3, 4}
	copy(y[:2], x)
	fmt.Println(y)

	switch "hello" {
	case "yes":
		fmt.Println("yes")
	default:
		fmt.Println("default")
	}

	c := make(chan bool)
	go func() {
		for d := range c {
			fmt.Println("range: ", d)
		}
	}()

	time.Sleep(time.Second * 1)
	c <- true
	c <- false
	c <- true
	close(c)
	time.Sleep(time.Second * 2)
}
