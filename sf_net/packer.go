package sf_net

import (
	"encoding/binary"
	"math"
	"test/sf_config"
	"test/sf_misc"
)

type Packer interface {
	pack(buf []byte) (Packet, error)
	unpack(Packet) []byte
	clear()
}

type PackerType string

const (
	PACKER_TYPE_DEFAULT PackerType = "default"
)

func New_packer(t PackerType) Packer {
	switch t {
	default:
		return new(packer)
	}

	return nil
}

func Get_packer_type(conf sf_config.Config) PackerType {
	t, e := conf.String()
	if e != nil {
		panic("stream type must be string")
	}
	switch t {
	default:
		return PACKER_TYPE_DEFAULT
	}
}

const (
	LEN_MSG_TYPE int = 4
	LEN_MSG_BODY int = 2
	LEN_MSG_TAIL int = 1
	LEN_MSG_HEAD int = LEN_MSG_TYPE + LEN_MSG_BODY + LEN_MSG_TAIL

	INDEX_MSG_TYPE int = 0
	INDEX_MSG_BODY int = INDEX_MSG_TYPE + LEN_MSG_TYPE
	INDEX_MSG_TAIL int = INDEX_MSG_BODY + LEN_MSG_BODY
)

type packer struct {
	buf []byte
}

func Append(src []byte, app []byte) (result []byte) {
	if cap(src)-len(src) < len(app) {
		result = make([]byte, len(src)+len(app))
		copy(result, src)
		result = result[:len(src)]
	} else {
		result = src
	}

	copy(result[len(result):], app)
	return result[:len(result)+len(app)]
}

func (p *packer) pack(buf []byte) (pkt Packet, err error) {
	if p == nil || buf == nil {
		err = sf_misc.ErrNilPointer
		return
	}

	if p.buf == nil {
		p.buf = make([]byte, len(buf))
	}

	p.buf = Append(p.buf, buf)
	if len(p.buf) < LEN_MSG_HEAD {
		return
	}

	var endian = binary.BigEndian
	t := int(int32(endian.Uint32(p.buf[INDEX_MSG_TYPE : INDEX_MSG_TYPE+LEN_MSG_TYPE])))
	body := int(int16(endian.Uint16(p.buf[INDEX_MSG_BODY : INDEX_MSG_BODY+LEN_MSG_BODY])))
	tail := int(int8(p.buf[INDEX_MSG_TAIL]))
	if 0 > body || 0 > tail {
		err = sf_misc.ErrPackerHeaderParse
		return
	}

	if len(p.buf) < LEN_MSG_HEAD+body+tail {
		return
	}

	data := new(packet)
	data.t = MsgType(t)
	data.body = make([]byte, body)
	copy(data.body, p.buf[LEN_MSG_HEAD:LEN_MSG_HEAD+body])
	data.tail = make([]byte, tail)
	copy(data.tail, p.buf[LEN_MSG_HEAD+body:LEN_MSG_HEAD+body+tail])

	copy(p.buf, p.buf[LEN_MSG_HEAD+body+tail:])
	p.buf = p.buf[:len(p.buf)-LEN_MSG_HEAD+body+tail]

	pkt = data
	return
}

func (p *packer) unpack(datap Packet) []byte {
	if p == nil || datap == nil {
		return nil
	}
	l := LEN_MSG_HEAD
	if datap.Get_tail() != nil {
		l += len(datap.Get_tail())
	}

	if datap.Get_body() != nil {
		l += len(datap.Get_body())
	}

	buf := make([]byte, l)

	endian := binary.BigEndian
	endian.PutUint32(buf[INDEX_MSG_TYPE:INDEX_MSG_TYPE+LEN_MSG_TYPE], uint32(datap.Get_type()))
	var body int = 0
	if datap.Get_body() != nil && len(datap.Get_body()) > math.MaxUint16 {
		panic("length of packet body bigger than math.MaxUint16")
	} else if datap.Get_body() != nil {
		body = len(datap.Get_body())
	}
	var tail int = 0
	if datap.Get_tail() != nil && len(datap.Get_tail()) > math.MaxUint8 {
		panic("length of packet tail bigger than math.MaxUint8")
	} else if datap.Get_tail() != nil {
		tail = len(datap.Get_tail())
	}

	endian.PutUint16(buf[INDEX_MSG_BODY:INDEX_MSG_BODY+LEN_MSG_BODY], uint16(body))
	buf[INDEX_MSG_TAIL] = byte(tail)
	copy(buf[LEN_MSG_HEAD:LEN_MSG_HEAD+body], datap.Get_body())
	copy(buf[LEN_MSG_HEAD+body:LEN_MSG_HEAD+body+tail], datap.Get_tail())

	return buf
}

func (p *packer) clear() {
	if p.buf != nil {
		p.buf = nil
	}
}
