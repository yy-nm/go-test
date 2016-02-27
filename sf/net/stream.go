package net

import (
	"encoding/binary"
	"io"
	"math"
	"test/sf/config"
	"test/sf/misc"
)

type StreamType string

const (
	STREAM_TYPE_DEFAULT StreamType = "default"
	STREAM_TYPE_PACKET  StreamType = "packet"
)

func newStream(t StreamType, src io.ReadWriteCloser) Stream {
	switch t {
	case STREAM_TYPE_PACKET:
		return &packetStream{stream{src, nil}}
	default:
		return &stream{conn: src}
	}

	return nil
}

func getStreamType(conf config.Config) StreamType {
	t, e := conf.String()
	if e != nil {
		panic("stream type must be string")
	}

	switch StreamType(t) {
	case STREAM_TYPE_PACKET:
		return STREAM_TYPE_PACKET
	default:
		return STREAM_TYPE_DEFAULT
	}
}

type Stream interface {
	Write(p Packet) (err error)
	Read() (p Packet, err error)
	Close() (err error)
}

type stream struct {
	conn io.ReadWriteCloser
	buf  []byte
}

func (s *stream) Close() error {
	if s != nil && s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *stream) Read() (p Packet, err error) {
	if s == nil || s.conn == nil {
		err = misc.ErrNilPointer
		return
	}

	if s.buf == nil {
		s.buf = make([]byte, 4*1024)
	}

	n, e := s.conn.Read(s.buf)
	if e != nil {
		err = e
		return
	}

	p = NewPacket(MSG_FLAG_REQ, s.buf[:n], nil)
	return
}

func (s *stream) Write(p Packet) (err error) {
	if s == nil || s.conn == nil {
		err = misc.ErrNilPointer
		return
	}

	buf := p.GetBody()
	if p == nil || buf == nil {
		return
	}

	var n int = 0
	for len(buf) > 0 {
		n, err = s.conn.Write(buf)
		if err != nil {
			return
		}
		buf = buf[n:]
	}
	return
}

const (
	LEN_MSG_TYPE int = 4
	LEN_MSG_BODY int = 2
	LEN_MSG_TAIL int = 1
	LEN_MSG_HEAD int = LEN_MSG_TYPE + LEN_MSG_BODY + LEN_MSG_TAIL

	INDEX_MSG_TYPE int = 0
	INDEX_MSG_BODY int = INDEX_MSG_TYPE + LEN_MSG_TYPE
	INDEX_MSG_TAIL int = INDEX_MSG_BODY + LEN_MSG_BODY

	LEN_MSG_BODY_MAX = math.MaxUint16
	LEN_MSG_TAIL_MAX = math.MaxUint8
)

type packetStream struct {
	stream
}

func (s *packetStream) Close() error {
	if s != nil {
		if s.buf != nil {
			s.buf = s.buf[0:0]
		}
		if s.conn != nil {
			return s.conn.Close()
		}
	}

	return nil
}

func readFromConn(s *packetStream) (err error) {
	if s.conn == nil {
		err = misc.ErrNilPointer
		return
	}
	len_old := len(s.buf)
	var n int = 0

	n, err = s.conn.Read(s.buf[len_old:cap(s.buf)])
	if err != nil {
		s.buf = s.buf[0:0]
	} else if n <= 0 {
		err = misc.ErrConnRead
	} else {
		s.buf = s.buf[:len_old+n]
	}

	return
}

func readFromBuf(s *packetStream) (p Packet, err error) {
	if s.buf == nil {
		s.buf = make([]byte, 4*1024)
		s.buf = s.buf[0:0]
	}

	for len(s.buf) < LEN_MSG_HEAD {
		err = readFromConn(s)
		if err != nil {
			return
		}
	}

	endian := binary.BigEndian
	t := MsgType(uint32(endian.Uint32(s.buf[INDEX_MSG_TYPE : INDEX_MSG_TYPE+LEN_MSG_TYPE])))
	body := int(uint16(endian.Uint16(s.buf[INDEX_MSG_BODY : INDEX_MSG_BODY+LEN_MSG_BODY])))
	tail := int(uint8(s.buf[INDEX_MSG_TAIL]))

	if body > LEN_MSG_BODY_MAX {
		err = misc.ErrPacketBodyLenTooBig
		return
	}

	if tail > LEN_MSG_TAIL_MAX {
		err = misc.ErrPacketTailLength
		return
	}

	if body+tail+LEN_MSG_HEAD > cap(s.buf) {
		buf := make([]byte, body+tail+LEN_MSG_HEAD)
		copy(buf, s.buf)
		buf = buf[:len(s.buf)]
		s.buf = buf
	}

	for body+tail+LEN_MSG_HEAD > len(s.buf) {
		err = readFromConn(s)
		if err != nil {
			return
		}
	}

	p = NewPacket(t, s.buf[LEN_MSG_HEAD:LEN_MSG_HEAD+body], s.buf[LEN_MSG_HEAD+body:LEN_MSG_HEAD+body+tail])

	copy(s.buf, s.buf[LEN_MSG_HEAD+body+tail:])
	s.buf = s.buf[:LEN_MSG_HEAD+body+tail]

	return
}

func (s *packetStream) Read() (p Packet, err error) {
	if s == nil || s.conn == nil {
		err = misc.ErrNilPointer
		return
	}

	return readFromBuf(s)
}

func (s *packetStream) Write(p Packet) (err error) {
	if s == nil || s.conn == nil {
		err = misc.ErrNilPointer
		return
	}

	l := LEN_MSG_HEAD
	if p.GetTail() != nil {
		l += len(p.GetTail())
	}

	if p.GetBody() != nil {
		l += len(p.GetBody())
	}

	buf := make([]byte, l)

	endian := binary.BigEndian
	copy(buf[INDEX_MSG_TYPE:INDEX_MSG_TYPE+LEN_MSG_TYPE], p.GetType().Convert2Bytes(endian))
	var body int = 0
	if p.GetBody() != nil && len(p.GetBody()) > LEN_MSG_BODY_MAX {
		panic("length of packet body bigger than LEN_MSG_BODY_MAX")
	} else if p.GetBody() != nil {
		body = len(p.GetBody())
	}
	var tail int = 0
	if p.GetTail() != nil && len(p.GetTail()) > LEN_MSG_TAIL_MAX {
		panic("length of packet tail bigger than LEN_MSG_TAIL_MAX")
	} else if p.GetTail() != nil {
		tail = len(p.GetTail())
	}

	endian.PutUint16(buf[INDEX_MSG_BODY:INDEX_MSG_BODY+LEN_MSG_BODY], uint16(body))
	buf[INDEX_MSG_TAIL] = byte(tail)
	copy(buf[LEN_MSG_HEAD:LEN_MSG_HEAD+body], p.GetBody())
	copy(buf[LEN_MSG_HEAD+body:LEN_MSG_HEAD+body+tail], p.GetTail())

	var n int = 0
	for len(buf) > 0 {
		n, err = s.conn.Write(buf)
		if err != nil {
			return
		}
		buf = buf[n:]
	}
	return
}
