package sf_net

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"

	"test/sf_config"
	"test/sf_misc"
)

var g_conn_id uint32 = 1

func get_new_id() uint32 {
	return atomic.AddUint32(&g_conn_id, 1)
}

func gen_uni_conn_id() conn_id {
	var i uint32
	for {
		i = get_new_id()
		if conn_id_invalid != conn_id(i) {
			break
		}
	}

	return conn_id(i)
}

type io_type struct {
	stype ServiceType
	sid   ServiceId
}

type io_svr struct {
	io_type
	l  *listener
	st StreamType
	pt PackerType
}

// TODO may use better struct or interface to instead io_client represent client and recver
type io_client struct {
	io_type
	c         *conn
	p         Packer
	s         Stream
	st        StreamType
	id        conn_id
	w_channel chan []byte
	is_recv   bool
}

type Network interface {
	send_msg(m Msg)
	close_conn(m Msg)

	register_callback(callback NetworkCallback) (old NetworkCallback)
	start() (err error)
	stop() (err error)
}

type NetworkCallback interface {
	recv_msg(m Msg)
	conn_broken(conn_id conn_id)
}

type conn_id int32

var conn_id_invalid conn_id = 0

// TODO network is too heavy too maintain, reduce
type network struct {
	svrs      list.List
	clients   map[ServiceType]map[ServiceId]*io_client
	conn_lock sync.Locker
	conns     map[conn_id]*io_client

	callback NetworkCallback
}

func conn_close(c *io_client) {
	if c.s != nil {
		c.s.Close()
	}
	if c.p != nil {
		c.p.clear()
	}
	if c.w_channel != nil {
		close(c.w_channel)
	}
}

func (n *network) handle_msg(m Msg) {
	if n.callback == nil {
		return
	}
	go n.callback.recv_msg(m)
}

func (n *network) handle_conn_close(c *io_client) {
	conn_close(c)
	n.conn_lock.Lock()
	delete(n.conns, c.id)
	n.conn_lock.Unlock()

	if n.callback == nil {
		return
	}

	go n.callback.conn_broken(c.id)
}

func (n *network) recv_data(c *io_client) {
	var p Packet
	buf := make([]byte, 4*1024)
	var l int
	var e error
	for {
		l, e = c.s.Read(buf)
		if e != nil {
			c.s.Close()
			break
		}
		if l > cap(buf) {
			buf = make([]byte, l)
			continue
		} else if l > len(buf) {
			buf = buf[:cap(buf)]
			continue
		}
		buf = buf[:l]
		p, e = c.p.pack(buf)
		if e != nil {
			c.s.Close()
			break
		}
		if p != nil {
			n.handle_msg(New_msg(p.Get_type(), p.Get_body(), p.Get_tail(), c.id))
		}
	}

	// when w_channel is close means client not longer reconn
	if !c.is_recv || c.w_channel != nil {
		n.client_conn(c)
	} else if c.is_recv {
		n.handle_conn_close(c)
	}
}

func (n *network) send_data(c *io_client) {
	var l int
	var e error
	for data := range c.w_channel {
		l = 0
		e = nil
		for len(data) > 0 {
			l, e = c.s.Write(data)
			if e != nil {
				c.s.Close() // case corresponding recv_data to read err
				return
			}
			data = data[l:]
		}
	}

	recover() // in case get panic
}

func (n *network) build_new_conn(c *io_client) {
	n.conn_lock.Lock()
	n.conns[c.id] = c
	n.conn_lock.Unlock()

	go n.recv_data(c)
	n.send_data(c)
}

// this func with endless-loop
func (n *network) recv_conn(s *io_svr) {
	for {
		c, err := s.l.recv()
		if err != nil {
			break
		}
		go n.build_new_conn(new_client(c, s)) // use gorountine  to avoid reduce accept speed
	}

	s.l.close()
	var recount int = 10
	for ; recount > 0; recount-- {
		err := s.l.connect()
		if err == nil {
			go n.recv_conn(s)
			return
		}
		s.l.close()
		time.Sleep(time.Second * 1)
	}

	// TODO the last op need
	panic("svr of transform fail")
}

// this func will block until client connect to target or err raise
func (n *network) client_conn(c *io_client) {
	if c.is_recv {
		//		panic(sf_misc.ErrNetworkRecvConnCannotConn.Error())
		return
	}

	if c.w_channel != nil {
		close(c.w_channel)
		c.w_channel = nil
	}
	for {
		err := c.c.connect()
		if err != nil {
			time.Sleep(time.Second * 2) // TODO retry issue, maybe use config in other type Network
			continue
		}
		break
	}
	c.w_channel = make(chan []byte, 1)
	c.s = New_Stream(c.st, c.c)
	if conn_id_invalid == c.id {
		c.id = gen_uni_conn_id()
		n.conn_lock.Lock()
		n.conns[c.id] = c
		n.conn_lock.Unlock()
	}

	go n.recv_data(c)
	go n.send_data(c)
}

// not thread-safe
func (n *network) start() (err error) {
	if n == nil {
		err = sf_misc.ErrNilPointer
		return
	}

	for e := n.svrs.Front(); e != nil; e = e.Next() {
		s, ok := e.Value.(*io_svr)
		if ok {
			err := s.l.connect()
			if err != nil {
				panic("start svr err, " + err.Error())
			}
		}
		go n.recv_conn(s)
	}

	for _, m := range n.clients {
		if m == nil {
			continue
		}
		for _, v := range m {
			v.id = conn_id_invalid
			go n.client_conn(v)
		}
	}

	return
}

// not thread-safe
func (n *network) stop() (err error) {
	for _, m := range n.clients {
		if m == nil {
			continue
		}
		for _, c := range m {
			if c == nil {
				continue
			}
			// avoid client invoke reconn
			close(c.w_channel)
			c.w_channel = nil
			c.p.clear()
			c.s.Close()
			n.conn_lock.Lock()
			delete(n.conns, c.id)
			n.conn_lock.Unlock()
			c.id = conn_id_invalid
		}
	}

	n.conn_lock.Lock()
	for k, c := range n.conns {
		delete(n.conns, k)
		if c == nil {
			continue
		}
		conn_close(c)
	}
	n.conn_lock.Unlock()

	for e := n.svrs.Front(); e != nil; e = e.Next() {
		s, ok := e.Value.(*io_svr)
		if ok {
			s.l.close()
		}
	}
	return
}

func (n *network) register_callback(callback NetworkCallback) (old NetworkCallback) {
	if n == nil {
		return
	}

	old, n.callback = n.callback, callback
	return
}

func (n *network) find_conn_by_msg(m Msg) *io_client {
	n.conn_lock.Lock()
	c, ok := n.conns[m.Get_id()]
	n.conn_lock.Unlock()

	if ok {
		return c
	} else {
		return nil
	}
}

func (n *network) close_conn(m Msg) {
	if n == nil || m == nil {
		return
	}
	c := n.find_conn_by_msg(m)
	conn_close(c)
}

func (n *network) send_msg(m Msg) {
	if n == nil || m == nil {
		return
	}
	c := n.find_conn_by_msg(m)
	if c.w_channel == nil {
		return
	}
	c.w_channel <- c.p.unpack(m)
}

// read config about "net" : {svrs, clients}
// demo config file
/*
"net" : {
	"type": "default"
	, "svrs": [{ "type": "gate", "id": 1, "sock_type": "tcp"
	, "sock_addr": "0.0.0.0:8080", "packer": "default", "stream": "default" }]
	, "clients": [{ "type": "gate", "id": 1, "sock_type": "tcp"
	, "sock_addr": "0.0.0.0:8081", "packer": "default", "stream": "default" }]
}
*/
const (
	CONF_NETWORK_KEY_TYPE = "type"

	CONF_NETWORK_KEY_SVR    = "svrs"
	CONF_NETWORK_KEY_CLIENT = "clients"

	CONF_NETWORK_KEY_SOCK_TYPE = "sock_type"
	CONF_NETWORK_KEY_SOCK_ADDR = "sock_addr"
	CONF_NETWORK_KEY_PACKER    = "packer"
	CONF_NETWORK_KEY_STREAM    = "stream"
)

func New_network(config sf_config.Config) (n Network, err error) {
	c, _ := config.Get(CONF_NETWORK_KEY_TYPE)
	s, _ := c.String()
	switch s {
	case "default":
		n, err = New_default_network(config)
	default:
		n, err = New_default_network(config)
	}

	return
}

func New_default_network(config sf_config.Config) (n Network, err error) {
	net := new(network)
	net.callback = nil
	//	net.ch_close = make(chan bool)
	net.conns = make(map[conn_id]*io_client)
	net.svrs = list.List{}
	net.svrs.Init()
	net.clients = make(map[ServiceType]map[ServiceId]*io_client)
	net.conn_lock = new(sync.Mutex)

	c, _ := config.Get(CONF_NETWORK_KEY_SVR)
	t, _ := c.Type()
	switch t {
	case sf_config.CONF_VAL_TYPE_ARR:
		for _, v := range c.Arr() {
			s := New_svr(v)
			net.svrs.PushBack(s)
		}
	case sf_config.CONF_VAL_TYPE_OBJ:
		s := New_svr(c)
		net.svrs.PushBack(s)
	default:
		panic("default network config err")
	}

	c, _ = config.Get(CONF_NETWORK_KEY_CLIENT)
	t, _ = c.Type()
	switch t {
	case sf_config.CONF_VAL_TYPE_ARR:
		for _, v := range c.Arr() {
			client := New_client(v)
			m, ok := net.clients[client.stype]
			if !ok {
				id_m := make(map[ServiceId]*io_client)
				net.clients[client.stype] = id_m
				m = id_m
			}
			m[client.sid] = client
		}
	case sf_config.CONF_VAL_TYPE_OBJ:
		client := New_client(c)
		m, ok := net.clients[client.stype]
		if !ok {
			id_m := make(map[ServiceId]*io_client)
			net.clients[client.stype] = id_m
			m = id_m
		}
		m[client.sid] = client
	}

	n = net
	return
}

const (
	CONF_NETWORK_SVR_KEY_TYPE = "type"
	CONF_NETWORK_SVR_KEY_ID   = "id"
)

func New_svr(config sf_config.Config) (s *io_svr) {
	s = new(io_svr)
	s.l = new(listener)

	c, _ := config.Get(CONF_NETWORK_SVR_KEY_TYPE)
	s.stype = Get_service_type(c)
	c, _ = config.Get(CONF_NETWORK_SVR_KEY_ID)
	s.sid = Get_service_id(c)

	c, _ = config.Get(CONF_NETWORK_KEY_PACKER)
	s.pt = Get_packer_type(c)
	c, _ = config.Get(CONF_NETWORK_KEY_STREAM)
	s.st = Get_stream_type(c)

	c, _ = config.Get(CONF_NETWORK_KEY_SOCK_TYPE)
	s.l.net_type = Config_get_string(c)
	c, _ = config.Get(CONF_NETWORK_KEY_SOCK_ADDR)
	s.l.net_addr = Config_get_string(c)

	return
}

func New_client(config sf_config.Config) (client *io_client) {
	client = new(io_client)
	client.c = new(conn)
	client.is_recv = false
	//	client.w_channel = make(chan []byte, 1)
	client.w_channel = nil // after conn or recv, or alway nil
	client.id = conn_id_invalid

	c, _ := config.Get(CONF_NETWORK_SVR_KEY_TYPE)
	client.stype = Get_service_type(c)
	c, _ = config.Get(CONF_NETWORK_SVR_KEY_ID)
	client.sid = Get_service_id(c)

	c, _ = config.Get(CONF_NETWORK_KEY_PACKER)
	client.p = New_packer(Get_packer_type(c))
	c, _ = config.Get(CONF_NETWORK_KEY_STREAM)
	client.st = Get_stream_type(c)

	c, _ = config.Get(CONF_NETWORK_KEY_SOCK_TYPE)
	client.c.net_type = Config_get_string(c)
	c, _ = config.Get(CONF_NETWORK_KEY_SOCK_ADDR)
	client.c.net_addr = Config_get_string(c)

	return
}

func new_client(conn *conn, s *io_svr) (c *io_client) {
	c = new(io_client)
	c.id = gen_uni_conn_id()
	c.c = conn
	c.is_recv = true
	c.p = New_packer(s.pt)
	c.s = New_Stream(s.st, conn)
	c.w_channel = make(chan []byte, 1)

	return
}
