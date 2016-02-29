package main

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"test/unit"
	"time"
)

type data struct {
	send int64
	recv int64
}

func (d *data) CollectSend(sendCount int) {
	d.send = atomic.AddInt64(&d.send, int64(sendCount))
}
func (d *data) CollectRecv(recvCount int) {
	d.recv = atomic.AddInt64(&d.recv, int64(recvCount))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("process need two param at least")
	}

	count, _ := strconv.Atoi(os.Args[1])
	port, ip, t, tm := unit.ParseData(os.Args[2:])

	d := new(data)
	var i int = 0
	for i < count {
		i++
		go unit.UnitTest(port, ip, t, tm, d)
	}

	time.Sleep(time.Second * time.Duration(tm))

	time.Sleep(time.Second)

	fmt.Println("total send speed: ", d.send)
	fmt.Println("total recv speed: ", d.recv)
}
