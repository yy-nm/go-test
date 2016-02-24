package main

import (
	"fmt"
	"test/sf_config"
	"test/sf_net"
	"time"
)

type svr struct {
	net sf_net.Network
}

func (s *svr) Recv_msg(m sf_net.Msg) {
	fmt.Print("recv: ")
	fmt.Println(string(m.Get_body()))
	s.net.Send_msg(m)
}

func (s *svr) Conn_broken(conn_id sf_net.Conn_Id) {
	fmt.Println("conn close: ", conn_id)
}

func main() {
	conf := `
{
	"net" : {
		"type": "default"
		, "svrs": [{ "type": "gate", "id": 1, "sock_type": "tcp"
		, "sock_addr": "0.0.0.0:9090", "packer": "default", "stream": "default" }]

	}
}
`
	j := sf_config.New_json_config()
	j.Read([]byte(conf))
	j, _ = j.Get("net")
	fmt.Println(j)

	n, e := sf_net.New_network(j)
	if e != nil {
		fmt.Println("err: ", e)
		return
	}

	s := new(svr)
	s.net = n
	n.Register_callback(s)
	n.Start()

	for {
		time.Sleep(time.Second * 2)
	}
}
