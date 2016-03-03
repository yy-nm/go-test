package net

import (
	"container/list"
	"net"
	"sync"
	"test/sf/config"
	"test/sf/misc"
	"time"
)

type ConnMgrCallback interface {
	ConnRecv(id ConnId, m Msg)
	ConnBroken(id ConnId)
	ConnBuild(id ConnId, addr *NetAddr)
}

type ConnMgr interface {
	Start() (err error)
	Stop() (err error)

	CloseConn(id ConnId)
	SendMsg(id ConnId, m Msg)

	Register(cb ConnMgrCallback) (old ConnMgrCallback)
}

type connState int

const (
	CONN_STATE_INIT connState = iota
	CONN_STATE_CONNECTING
	CONN_STATE_CONNECTED
	CONN_STATE_BROKEN
)

type State int

const (
	STATE_INIT State = iota
	STATE_START
	STATE_STOP
)

type conn struct {
	//clt   IOClt
	IOClt
	state connState
}

type connManager struct {
	conn_lock sync.Mutex
	conns     map[ConnId]*conn
	svrs      *list.List
	clts      *list.List
	callback  ConnMgrCallback
	state     State
}

func (cm *connManager) CloseConn(id ConnId) {
	if cm == nil {
		return
	}

	cm.conn_lock.Lock()
	defer cm.conn_lock.Unlock()

	c, ok := cm.conns[id]
	if !ok { // client conn can not be closed
		return
	}
	delete(cm.conns, id)
	c.state = CONN_STATE_BROKEN
	c.Close()
	if cm.callback != nil {
		go cm.callback.ConnBroken(id)
	}
}

func (cm *connManager) Register(cb ConnMgrCallback) (old ConnMgrCallback) {
	if cm == nil {
		return
	}

	old, cm.callback = cm.callback, cb
	return
}

func (cm *connManager) Start() (err error) {
	if cm == nil {
		err = misc.ErrNilPointer
		return
	}

	cm.state = STATE_START
	for e := cm.svrs.Front(); e != nil; e = e.Next() {
		if e.Value == nil {
			continue
		}

		svr, ok := e.Value.(IOSvr)
		if !ok {
			panic("connManager svrs contains non-IOSvr elements")
		}
		err = svr.Connect()
		if err != nil {
			return
		}

		go cm.recvClients(svr)
	}

	for e := cm.clts.Front(); e != nil; e = e.Next() {
		if e.Value == nil {
			continue
		}
		client, ok := e.Value.(*conn)
		if !ok {
			panic("connManager clts contains non-conn elements")
		}
		go cm.clientConn(client)
	}

	return
}

func (cm *connManager) Stop() (err error) {
	if cm == nil {
		err = misc.ErrNilPointer
		return
	}

	cm.state = STATE_STOP
	for e := cm.svrs.Front(); e != nil; e = e.Next() {
		if e.Value == nil {
			continue
		}

		svr, ok := e.Value.(IOSvr)
		if !ok {
			panic("connManager svrs contains non-IOSvr elements")
		}
		svr.Close()
	}

	for e := cm.clts.Front(); e != nil; e = e.Next() {
		if e.Value == nil {
			continue
		}
		client, ok := e.Value.(*conn)
		if !ok {
			panic("connManager clts contains non-conn elements")
		}
		client.state = CONN_STATE_BROKEN
		client.Close()
	}

	for k, _ := range cm.conns {
		cm.CloseConn(k)
	}

	return
}

func (cm *connManager) SendMsg(id ConnId, m Msg) {
	if cm == nil || cm.state != STATE_START {
		return
	}
	cm.conn_lock.Lock()
	c, ok := cm.conns[id]
	cm.conn_lock.Unlock()
	if !ok {
		return
	}

	if c.state == CONN_STATE_BROKEN || c.state == CONN_STATE_INIT {
		go cm.clientReConn(c)
	} else if c.state == CONN_STATE_CONNECTED {
		go cm.clientSendMsg(c, m)
	}
}

func (cm *connManager) recvClients(s IOSvr) {
	if cm == nil || s == nil {
		return
	}

	for {
		n, err := s.Accept()
		if err != nil {
			break
		}
		go cm.buildNewClients(n, s)
	}

	s.Close()
	var recount int = 10
	for ; recount > 0 && cm.state == STATE_START; recount-- {
		err := s.Connect()
		if err == nil {
			go cm.recvClients(s)
			return
		}
		s.Close()
		time.Sleep(time.Second * 1)
	}

	if cm.state == STATE_STOP {
		return
	}

	panic("Svr listen err after several retry")
}

func (cm *connManager) buildNewClients(n net.Conn, s IOSvr) {
	st, e := s.StreamType()
	if e != nil {
		return
	}

	client := newSvrClient(n, st, GenUniConnId())
	c := &conn{IOClt: client, state: CONN_STATE_CONNECTED}
	cm.addNewClients(c)
}

func (cm *connManager) addNewClients(c *conn) {
	if cm == nil || c == nil {
		return
	}

	cm.conn_lock.Lock()
	cm.conns[c.GetId()] = c
	cm.conn_lock.Unlock()

	go cm.recvMsg(c)

	if cm.callback != nil {
		cm.callback.ConnBuild(c.GetId(), nil)
	}
}

func (cm *connManager) clientConn(c *conn) {
	if cm == nil || c == nil {
		return
	}

	c.Close()
	c.state = CONN_STATE_INIT
	for {
		err := cm.clientReConn(c)
		if err == nil {
			break
		}
	}
}

func (cm *connManager) clientReConn(c *conn) (err error) {
	if cm == nil || c == nil {
		return
	}

	if cm.state == STATE_STOP {
		err = misc.ErrConManagerStop
		return
	}
	if c.state == CONN_STATE_CONNECTED || c.state == CONN_STATE_CONNECTING {
		return
	}

	c.state = CONN_STATE_CONNECTING
	err = c.Connect()
	if err != nil {
		if c.state == CONN_STATE_CONNECTING {
			c.state = CONN_STATE_BROKEN
		}
		return
	}

	c.state = CONN_STATE_CONNECTED
	cm.addNewClients(c)
	return
}

func (cm *connManager) recvMsg(c *conn) {
	if cm == nil || c == nil {
		return
	}

	for {
		if c.state != CONN_STATE_CONNECTED {
			break
		}
		p, e := c.Recv()
		if e != nil {
			c.state = CONN_STATE_BROKEN
			break
		}

		//misc.Log("read from client: ", string(p.GetBody()), ", id: ", c.clt.GetId())
		if p != nil && cm.callback != nil {
			go cm.callback.ConnRecv(c.GetId(), NewMsg(p.GetType(), p.GetBody(), p.GetTail(), c.GetId()))
		}
	}

	if c.state == CONN_STATE_BROKEN && c.IsRecv() {
		cm.CloseConn(c.GetId())
	}
}

func (cm *connManager) clientSendMsg(c *conn, m Msg) {
	if cm == nil || c == nil {
		return
	}

	if c.state == CONN_STATE_BROKEN && !c.IsRecv() {
		cm.clientReConn(c)
	} else if c.state == CONN_STATE_BROKEN {
		return
	}

	if c.state == CONN_STATE_CONNECTED {
		c.Send(m)
	}
}

const (
	CONF_CONNMGR_TYPE         string = "type"
	CONF_CONNMGR_TYPE_DEFAULT string = "default"

	CONF_NETWORK_KEY_SVR    = "svrs"
	CONF_NETWORK_KEY_CLIENT = "clients"

	CONF_NETWORK_KEY_SOCK_TYPE = "sock_type"
	CONF_NETWORK_KEY_SOCK_ADDR = "sock_addr"
	CONF_NETWORK_KEY_STREAM    = "stream"
)

func NewConnMgr(conf config.Config) ConnMgr {
	t, _ := conf.Get(CONF_CONNMGR_TYPE)

	switch v, _ := t.String(); v {
	case CONF_CONNMGR_TYPE_DEFAULT:
		fallthrough
	default:
		return newDefaultConnMgr(conf)
	}
}

func newDefaultConnMgr(conf config.Config) ConnMgr {
	cm := new(connManager)
	cm.conns = make(map[ConnId]*conn)
	cm.clts = new(list.List).Init()
	cm.svrs = new(list.List).Init()
	cm.state = STATE_INIT

	c, _ := conf.Get(CONF_NETWORK_KEY_SVR)
	t, _ := c.Type()
	switch t {
	case config.CONF_VAL_TYPE_ARR:
		for _, v := range c.Arr() {
			s := readSvr(v)
			cm.svrs.PushBack(s)
		}
	case config.CONF_VAL_TYPE_OBJ:
		s := readSvr(c)
		cm.svrs.PushBack(s)
	default:
		panic("default network config err")
	}

	c, _ = conf.Get(CONF_NETWORK_KEY_CLIENT)
	t, _ = c.Type()
	switch t {
	case config.CONF_VAL_TYPE_ARR:
		for _, v := range c.Arr() {
			client := readClient(v)
			cm.clts.PushBack(client)
		}
	case config.CONF_VAL_TYPE_OBJ:
		client := readClient(c)
		cm.clts.PushBack(client)
	}

	return cm
}

func readSvr(conf config.Config) IOSvr {
	c, _ := conf.Get(CONF_NETWORK_KEY_SOCK_TYPE)
	nt := ConfigGetString(c)
	c, _ = conf.Get(CONF_NETWORK_KEY_SOCK_ADDR)
	na := ConfigGetString(c)
	c, _ = conf.Get(CONF_NETWORK_KEY_STREAM)
	st := getStreamType(c)

	return newSvr(nt, na, st)
}

func readClient(conf config.Config) *conn {
	c := readClt(conf)
	if c == nil {
		return nil
	}

	return &conn{IOClt: c, state: CONN_STATE_INIT}
}

func readClt(conf config.Config) IOClt {
	c, _ := conf.Get(CONF_NETWORK_KEY_SOCK_TYPE)
	nt := ConfigGetString(c)
	c, _ = conf.Get(CONF_NETWORK_KEY_SOCK_ADDR)
	na := ConfigGetString(c)
	c, _ = conf.Get(CONF_NETWORK_KEY_STREAM)
	st := getStreamType(c)

	id := GenUniConnId()

	return newClient(nt, na, st, id)
}
