package net

import "encoding/binary"

type MsgType uint32

// MsgType
// 1111 xxxx xxxx xxxx xxxx xxxx xxxx xxxx: bit represent send/recv state
// xxxx 1111 1111 xxxx xxxx xxxx xxxx xxxx: opt, used for cluster, means cluster type name
// xxxx xxxx xxxx 1111 xxxx xxxx xxxx xxxx: num represent msg type(i.e. interal msg, system msg, normal msg)
// xxxx xxxx xxxx xxxx 1111 1111 1111 1111: num represent msg proto num
const (
	MSG_FLAG_NULL MsgType = 0x00000000
	MSG_FLAG_FULL MsgType = 0xffffffff

	MSG_FLAG_REQ           MsgType = 0x10000000
	MSG_FLAG_REQ_NOT       MsgType = ^MSG_FLAG_REQ
	MSG_FLAG_RES           MsgType = 0x20000000
	MSG_FLAG_RES_NOT       MsgType = ^MSG_FLAG_RES
	MSG_FLAG_UNICAST       MsgType = 0x40000000
	MSG_FLAG_UNICAST_NOT   MsgType = ^MSG_FLAG_UNICAST
	MSG_FLAG_MULTICAST     MsgType = 0x80000000
	MSG_FLAG_MULTICAST_NOT MsgType = ^MSG_FLAG_MULTICAST
	MSG_FLAG_SEND_STATE    MsgType = 0xf0000000

	MSG_FLAG_OPT MsgType = 0x0ff00000

	MSG_FLAG_TYPE     MsgType = 0x000f0000
	MSG_FLAG_TYPE_NOR MsgType = 0x00000000
	MSG_FLAG_TYPE_SYS MsgType = 0x00010000
	MSG_FLAG_TYPE_INT MsgType = 0x00020000

	MSG_FLAG_PROTO MsgType = 0x0000ffff

	MSG_SHIFT_OPT uint = 5
)

type Packet interface {
	GetType() MsgType
	GetBody() []byte
	GetTail() []byte
}

type Msg interface {
	Packet
	GetId() ConnId
}

type packet struct {
	t    MsgType
	body []byte
	tail []byte
}

type msg struct {
	packet
	id ConnId
}

func NewMsg(t MsgType, body []byte, tail []byte, id ConnId) Msg {
	m := new(msg)
	m.t = t
	if body != nil {
		m.body = body
	}
	if tail != nil {
		m.tail = tail
	}
	m.id = id
	return m
}

func NewPacket(t MsgType, body []byte, tail []byte) Packet {
	p := new(packet)
	p.t = t
	if body != nil {
		p.body = make([]byte, len(body))
		copy(p.body, body)
	}
	if tail != nil {
		p.tail = make([]byte, len(tail))
		copy(p.tail, tail)
	}

	return p
}

func (mt MsgType) IsRequest() bool {
	return mt&MSG_FLAG_REQ != MSG_FLAG_NULL
}

func (mt MsgType) IsResponse() bool {
	return mt&MSG_FLAG_RES != MSG_FLAG_NULL
}

func (mt MsgType) IsCluster() bool {
	return mt&MSG_FLAG_MULTICAST != MSG_FLAG_NULL || mt&MSG_FLAG_UNICAST != MSG_FLAG_NULL
}

func (mt MsgType) GetOpt() int {
	opt := mt & MSG_FLAG_OPT
	o := int(opt)
	o >>= MSG_SHIFT_OPT
	return o
}

func (mt MsgType) GetType() MsgType {
	return mt & MSG_FLAG_TYPE
}

func (mt MsgType) GetProto() int {
	return int(mt & MSG_FLAG_PROTO)
}

func (mt MsgType) SetRequestFlag(set_or_cancel bool) MsgType {
	if set_or_cancel {
		mt |= MSG_FLAG_REQ
	} else {
		mt &= MSG_FLAG_REQ_NOT
	}

	return mt
}

func (mt MsgType) SetResponseFlag(set_or_cancel bool) MsgType {
	if set_or_cancel {
		mt |= MSG_FLAG_RES
	} else {
		mt &= MSG_FLAG_RES_NOT
	}

	return mt
}

func (mt MsgType) SetUnicastFlag(set_or_cancel bool) MsgType {
	if set_or_cancel {
		mt |= MSG_FLAG_UNICAST
	} else {
		mt &= MSG_FLAG_UNICAST_NOT
	}

	return mt
}

func (mt MsgType) SetMulticastFlag(set_or_cancel bool) MsgType {
	if set_or_cancel {
		mt |= MSG_FLAG_MULTICAST
	} else {
		mt &= MSG_FLAG_MULTICAST_NOT
	}

	return mt
}

func (mt MsgType) SetOpt(opt int) MsgType {
	o := MsgType(opt)
	o <<= MSG_SHIFT_OPT
	o &= MSG_FLAG_OPT
	mt |= o

	return mt
}

func (mt MsgType) SetType(t MsgType) MsgType {
	t &= MSG_FLAG_TYPE
	mt |= t

	return mt
}

func (mt MsgType) SetProto(proto int) MsgType {
	p := MsgType(proto)
	p &= MSG_FLAG_PROTO
	mt |= p

	return mt
}

func (mt MsgType) Convert2Bytes(bo binary.ByteOrder) (buf []byte) {
	buf = make([]byte, 4)
	bo.PutUint32(buf, uint32(mt))
	return
}

func (p *packet) GetType() MsgType {
	return p.t
}

func (p *packet) GetBody() []byte {
	return p.body
}

func (p *packet) GetTail() []byte {
	return p.tail
}

func (m *msg) GetId() ConnId {
	return m.id
}
