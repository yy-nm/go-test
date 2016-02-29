package unit

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type Statistics interface {
	CollectSend(sendCount int)
	CollectRecv(recvCount int)
}

type DeviceZero struct {
	Count int
}

func (dz *DeviceZero) Read(p []byte) (n int, err error) {

	for i, _ := range p {
		p[i] = 0
	}

	n = len(p)
	dz.Count += n
	return
}

func (dz *DeviceZero) Close() error {
	return nil
}

type DeviceNull struct {
	Count int
}

func (dn *DeviceNull) Write(p []byte) (n int, err error) {
	n = len(p)
	dn.Count += n
	return
}

func (dn *DeviceNull) Close() error {
	return nil
}

func ParseData(args []string) (port int, ip, t string, tm int) {
	port = 9090
	ip = "127.0.0.1"
	t = "tcp"
	tm = 10

	switch len(args) {
	case 4:
		tm, _ = strconv.Atoi(args[3])
		fallthrough
	case 3:
		t = args[2]
		fallthrough
	case 2:
		ip = args[1]
		fallthrough
	case 1:
		port, _ = strconv.Atoi(args[0])
	}

	return
}

func UnitTest(port int, ip, t string, tm int, ss Statistics) {
	s, err := net.Dial(t, fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}

	dz := new(DeviceZero)
	dn := new(DeviceNull)

	go func() {
		io.Copy(s, dz)
	}()

	go func() {
		io.Copy(dn, s)
	}()

	time.Sleep(time.Second * time.Duration(tm))

	fmt.Printf("send speed: %d\n", dz.Count/tm)
	fmt.Printf("recv speed: %d\n", dn.Count/tm)

	if ss != nil {
		ss.CollectSend(dz.Count / tm)
		ss.CollectRecv(dn.Count / tm)
	}
	s.Close()
}
