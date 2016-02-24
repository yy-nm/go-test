// author mard

package sf_net

import (
	"net"
	"test/sf_misc"
)

type conn struct {
	c net.Conn

	net_type string
	net_addr string
}

func (c *conn) connect() (err error) {
	if c == nil {
		return sf_misc.ErrNilPointer
	}

	if c.c != nil {
		c.c.Close()
		c.c = nil
	}

	c.c, err = net.Dial(c.net_type, c.net_addr)

	return
}

func (c *conn) Read(buf []byte) (n int, err error) {
	if c == nil || c.c == nil {
		n = 0
		err = sf_misc.ErrNilPointer
		return
	}

	n, err = c.c.Read(buf)
	return
}

func (c *conn) Write(buf []byte) (n int, err error) {
	if c == nil || c.c == nil {
		n = 0
		err = sf_misc.ErrNilPointer
		return
	}

	var l int = 0
	n = 0
	for len(buf) > 0 {
		l, err = c.c.Write(buf[l:])
		if err != nil {
			return
		}
		n += l
		buf = buf[l:]
		l = 0
	}

	return
}

func (c *conn) Close() (err error) {
	if c == nil {
		return
	}

	if c.c != nil {
		err = c.c.Close()
		c.c = nil
	}
	return
}

type listener struct {
	l        net.Listener
	net_type string
	net_addr string
}

func (l *listener) connect() (err error) {
	if l == nil {
		err = sf_misc.ErrNilPointer
		return
	}
	if l.l != nil {
		l.l.Close()
		l.l = nil
	}

	l.l, err = net.Listen(l.net_type, l.net_addr)

	return
}

func (l *listener) recv() (c *conn, err error) {
	if l == nil || l.l == nil {
		err = sf_misc.ErrNilPointer
		return
	}

	nc, e := l.l.Accept()
	if e != nil {
		err = e
		return
	}

	c = New_conn_by_net(nc)
	return
}

func (l *listener) close() {
	if l == nil {
		return
	}

	if l.l != nil {
		l.l.Close()
		l.l = nil
	}
}

func New_listener(t, addr string) *listener {

	return &listener{net_type: t, net_addr: addr}
}

func New_conn(t, addr string) *conn {
	return &conn{net_type: t, net_addr: addr}
}

func New_conn_by_net(c net.Conn) *conn {
	return &conn{c: c}
}
