package net

import (
	"net"
	"test/sf/misc"
)

type IOClt interface {
	Connect() (err error)

	GetNetAddr() (na *NetAddr)

	Send(p Packet) (err error)
	Recv() (p Packet, err error)

	IsRecv() bool

	Close() (err error)

	GetId() ConnId
}

type IOSvr interface {
	Connect() (err error)
	Accept() (n net.Conn, err error)
	Close() (err error)
	StreamType() (st StreamType, err error)
}

type NetAddr struct {
	Type string
	Addr string
}

// recv from IO_svr
type svrClient struct {
	s  Stream
	id ConnId
}

func (sc *svrClient) Connect() (err error) {
	err = misc.ErrCanNotConn
	return
}

func (sc *svrClient) GetNetAddr() (na *NetAddr) {
	return
}

func (sc *svrClient) Send(p Packet) (err error) {
	if sc == nil || sc.s == nil {
		err = misc.ErrNilPointer
		return
	}

	err = sc.s.Write(p)
	return
}

func (sc *svrClient) Recv() (p Packet, err error) {
	if sc == nil || sc.s == nil {
		err = misc.ErrNilPointer
		return
	}

	p, err = sc.s.Read()
	return
}

func (sc *svrClient) IsRecv() bool {
	return true
}

func (sc *svrClient) Close() (err error) {
	if sc == nil || sc.s == nil {
		err = misc.ErrNilPointer
		return
	}

	err = sc.s.Close()
	return
}

func (sc *svrClient) GetId() ConnId {
	if sc == nil {
		return GetInvalidConnId()
	}

	return sc.id
}

type client struct {
	NetAddr
	svrClient
	st StreamType
}

func (c *client) Connect() (err error) {
	if c == nil {
		err = misc.ErrNilPointer
		return
	}

	var n net.Conn
	n, err = net.Dial(c.Type, c.Addr)
	if err != nil {
		return
	}

	c.s = newStream(c.st, n)
	return
}

func (c *client) IsRecv() bool {
	return false
}

func (c *client) GetNetAddr() (na *NetAddr) {
	if c == nil {
		return nil
	}
	return &c.NetAddr
}

type svr struct {
	NetAddr
	l  net.Listener
	st StreamType
}

func (s *svr) Connect() (err error) {
	if s == nil {
		err = misc.ErrNilPointer
		return
	}
	if s.l != nil {
		s.l.Close()
		s.l = nil
	}

	s.l, err = net.Listen(s.Type, s.Addr)
	return
}

func (s *svr) Accept() (n net.Conn, err error) {
	if s == nil || s.l == nil {
		err = misc.ErrNilPointer
		return
	}

	n, err = s.l.Accept()
	if err != nil {
		return
	}

	return
}

func (s *svr) Close() (err error) {
	if s == nil || s.l == nil {
		err = misc.ErrNilPointer
		return
	}

	err = s.l.Close()
	s.l = nil
	return
}

func (s *svr) StreamType() (st StreamType, err error) {
	if s == nil || s.l == nil {
		err = misc.ErrNilPointer
		return
	}

	st = s.st
	return
}

func newSvrClient(n net.Conn, st StreamType, id ConnId) IOClt {
	return &svrClient{s: newStream(st, n), id: id}
}

func newClient(t, addr string, st StreamType, id ConnId) IOClt {
	return &client{NetAddr{t, addr}, svrClient{nil, id}, st}
}

func newSvr(t, addr string, st StreamType) IOSvr {
	return &svr{NetAddr{t, addr}, nil, st}
}
