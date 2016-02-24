package sf_net

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
	Get_type() MsgType
	Get_body() []byte
	Get_tail() []byte
}

type Msg interface {
	Packet
	Get_id() conn_id
}

type packet struct {
	t    MsgType
	body []byte
	tail []byte
}

type msg struct {
	packet
	id conn_id
}

func New_msg(t MsgType, body []byte, tail []byte, id conn_id) Msg {
	m := new(msg)
	m.t = t
	m.body = body
	m.tail = tail
	m.id = id
	return m
}

func (mt MsgType) Is_flag_request() bool {
	return mt&MSG_FLAG_REQ != MSG_FLAG_NULL
}

func (mt MsgType) Is_flag_response() bool {
	return mt&MSG_FLAG_RES != MSG_FLAG_NULL
}

func (mt MsgType) Is_flag_cluster() bool {
	return mt&MSG_FLAG_MULTICAST != MSG_FLAG_NULL || mt&MSG_FLAG_UNICAST != MSG_FLAG_NULL
}

func (mt MsgType) Get_flag_opt() int {
	opt := mt & MSG_FLAG_OPT
	o := int(opt)
	o >>= MSG_SHIFT_OPT
	return o
}

func (mt MsgType) Get_flag_type() MsgType {
	return mt & MSG_FLAG_TYPE
}

func (mt MsgType) Get_flag_proto() MsgType {
	return mt & MSG_FLAG_PROTO
}

func (mt MsgType) Set_flag_request(set_or_cancel bool) MsgType {
	if set_or_cancel {
		mt |= MSG_FLAG_REQ
	} else {
		mt &= MSG_FLAG_REQ_NOT
	}

	return mt
}

func (mt MsgType) Set_flag_response(set_or_cancel bool) MsgType {
	if set_or_cancel {
		mt |= MSG_FLAG_RES
	} else {
		mt &= MSG_FLAG_RES_NOT
	}

	return mt
}

func (mt MsgType) Set_flag_unicast(set_or_cancel bool) MsgType {
	if set_or_cancel {
		mt |= MSG_FLAG_UNICAST
	} else {
		mt &= MSG_FLAG_UNICAST_NOT
	}

	return mt
}

func (mt MsgType) Set_flag_multicast(set_or_cancel bool) MsgType {
	if set_or_cancel {
		mt |= MSG_FLAG_MULTICAST
	} else {
		mt &= MSG_FLAG_MULTICAST_NOT
	}

	return mt
}

func (mt MsgType) Set_flag_opt(opt int) MsgType {
	o := MsgType(opt)
	o <<= MSG_SHIFT_OPT
	o &= MSG_FLAG_OPT
	mt |= o

	return mt
}

func (mt MsgType) Set_flag_type(t MsgType) MsgType {
	t &= MSG_FLAG_TYPE
	mt |= t

	return mt
}

func (mt MsgType) Set_flag_proto(proto int) MsgType {
	p := MsgType(proto)
	p &= MSG_FLAG_PROTO
	mt |= p

	return mt
}

func (p *packet) Get_type() MsgType {
	return p.t
}

func (p *packet) Get_body() []byte {
	return p.body
}

func (p *packet) Get_tail() []byte {
	return p.tail
}

func (m *msg) Get_id() conn_id {
	return m.id
}
