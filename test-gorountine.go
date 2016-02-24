package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	go func() {
		alarm := time.Tick(10 * time.Second)
		go func() {
			for {
				alarm := time.Tick(time.Second)
				fmt.Println(<-alarm)
			}
		}()
		<-alarm
	}()

	time.Sleep(time.Second * 20)
	net.TCPConn{}
}
