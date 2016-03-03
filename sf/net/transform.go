package net

import "container/list"

// used to get msg target id of conn, or list of target id
type Transform interface {
	RegisterService(s *Service, addr *NetAddr)
	RegisterConn(addr *Network, id ConnId)
	Transform(m *Msg) *list.List
}
