package main

import (
	"fmt"
	"runtime"
	"test/sf/config"
	"test/sf/net"
	"time"
)

type svr struct {
	mgr net.ConnMgr
}

func (s *svr) ConnRecv(id net.ConnId, m net.Msg) {
	//fmt.Print("recv: ")
	//fmt.Println(string(m.GetBody()))
	s.mgr.SendMsg(id, m)
}

func (s *svr) ConnBroken(id net.ConnId) {
	fmt.Println("conn close: ", id)
}

func (s *svr) ConnBuild(id net.ConnId, addr *net.NetAddr) {
	fmt.Println("new conn coming: id: ", id, ", addr: ", addr)
}

func main() {
	conf := `
{
	"net" : {
		"type": "default"
		, "conn" : {
			"svrs": [{ "type": "gate", "id": 1, "sock_type": "tcp"
			, "sock_addr": "0.0.0.0:9091", "stream": "default" }]
		}
	}
}
`
	j := config.NewJsonConfig()
	j.Read([]byte(conf))
	j, _ = j.Get("net")
	j, _ = j.Get("conn")
	fmt.Println(j)

	conns := net.NewConnMgr(j)

	s := new(svr)
	s.mgr = conns
	conns.Register(s)
	conns.Start()

	var ms runtime.MemStats
	for {
		time.Sleep(time.Second * 5)
		runtime.ReadMemStats(&ms)

		fmt.Printf("runtime: alloc mem: %d, free mem: %d, sys mem: %d\n", ms.Alloc, ms.Frees, ms.Sys)
		fmt.Printf("heap: alloc: %d, sys: %d, in-use: %d, no-use: %d\n", ms.HeapAlloc, ms.HeapSys, ms.HeapInuse, ms.HeapIdle)
	}
}
