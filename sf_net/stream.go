package sf_net

import (
	"encoding/binary"
	"fmt"
	"io"
	"test/sf_config"
	"test/sf_misc"
)

type StreamType string

const (
	STREAM_TYPE_DEFAULT StreamType = "default"
	STREAM_TYPE_PACKET  StreamType = "packet"
)

func New_Stream(t StreamType, src io.ReadWriteCloser) Stream {
	switch t {
	case STREAM_TYPE_PACKET:
		fmt.Println("packet")
		return &packet_stream{stream{src}, nil}
	default:
		fmt.Println("default")
		return &stream{conn: src}
	}

	return nil
}

func Get_stream_type(conf sf_config.Config) StreamType {
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
	io.ReadWriteCloser
}

type stream struct {
	conn io.ReadWriteCloser
}

func (s *stream) Close() error {
	if s != nil && s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *stream) Read(buf []byte) (n int, err error) {
	if s == nil || s.conn == nil {
		return 0, sf_misc.ErrNilPointer
	}

	return s.conn.Read(buf)
}

func (s *stream) Write(buf []byte) (n int, err error) {
	if s == nil || s.conn == nil {
		return 0, sf_misc.ErrNilPointer
	}

	return s.conn.Write(buf)
}

const (
	LEN_PACKET_HEAD int = 4
	LEN_PACKET_MAX  int = 4 * 1024 // max packet len 4k
)

type packet_stream struct {
	stream
	buf []byte
}

func (s *packet_stream) Close() error {
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

func read_from_conn(s *packet_stream) (n int, err error) {
	if s.conn == nil {
		return 0, sf_misc.ErrNilPointer
	}
	len_old := len(s.buf)

	n, err = s.conn.Read(s.buf[len_old:cap(s.buf)])
	if err != nil {
		s.buf = s.buf[0:0]
	} else if n < 0 {
		err = sf_misc.ErrConnRead
	} else {
		s.buf = s.buf[:len_old+n]
	}

	return
}

func read_from_buf(s *packet_stream, buf []byte) (n int, err error) {
	if s.buf == nil {
		s.buf = make([]byte, LEN_PACKET_MAX+LEN_PACKET_HEAD)
		s.buf = s.buf[0:0]
	}

	for len(s.buf) < LEN_PACKET_HEAD {
		_, err := read_from_conn(s)
		if err != nil {
			return 0, err
		}
	}

	endian := binary.BigEndian
	n = int(int32(endian.Uint32(s.buf[:LEN_PACKET_HEAD])))
	if n > LEN_PACKET_MAX {
		err = sf_misc.ErrPacketLengthTooBig
		return
	} else if 0 > n {
		err = sf_misc.ErrPacketLength
		return
	}

	for n+LEN_PACKET_HEAD > len(s.buf) {
		_, err = read_from_conn(s)
		if err != nil {
			return
		}
	}

	if n > len(buf) {
		return
	}

	copy(buf, s.buf[LEN_PACKET_HEAD:LEN_PACKET_HEAD+n])

	copy(s.buf, s.buf[n:])
	s.buf = s.buf[:len(s.buf)-LEN_PACKET_HEAD+n]

	return
}

// read **one** package at once dur call,
// if n = 0, means cannot read
// if n > len(buf), means length of buf is too short
func (s *packet_stream) Read(buf []byte) (n int, err error) {
	if s == nil || s.conn == nil {
		return 0, sf_misc.ErrNilPointer
	}

	return read_from_buf(s, buf)
}
