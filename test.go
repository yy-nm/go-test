package main

import (
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net"
	"test/sf_crypto"
	"time"
)

type Msg_package interface {
	get_type() int
	get_body() []byte
	get_tail() int
}

type Encode interface {
	unpack(*Msg_package) []byte
	pack([]byte) *Msg_package
}

type netio struct {
	Encode Encode //

}

type encode struct {
}

func (e *encode) pack(buf []byte) *Msg_package {
	if buf == nil || 0 == len(buf) {
		fmt.Println("buf is nil or len 0")
	}

	fmt.Printf("% x\n", buf)

	return nil
}

func (e *encode) unpack(m *Msg_package) (buf []byte) {
	if m == nil {
		fmt.Println("msg is nil")
	}

	fmt.Printf("msg type: %d\n", (*m).get_type())
	fmt.Printf("msg tail: %d\n", (*m).get_tail())

	return
}

func main() {

	ch_signal := make(chan int)
	ch_conn := make(chan net.Conn)

	for i := 0; i < 5; i++ {
		go func(i int) {
			xx := <-ch_signal
			fmt.Printf("index: %d, signal: %d", i, xx)
		}(i)
	}

	ch_signal <- 1
	ch_signal <- 2

	fmt.Println()
	rand.Seed(time.Now().UnixNano())
	key1 := big.NewInt(rand.Int63())
	key2 := big.NewInt(rand.Int63())

	fmt.Printf("key1: %d, key2: %d\n", key1, key2)

	key1_ := sf_crypto.DH_exchange(key1)
	key2_ := sf_crypto.DH_exchange(key2)

	fmt.Printf("key1': %d, key2': %d\n", key1_, key2_)

	fmt.Println(sf_crypto.DH_secret(key2_, key1))
	fmt.Println(sf_crypto.DH_secret(key1_, key2))

	key := sf_crypto.DH_secret(key2_, key1)
	var b = key.Bytes()
	fmt.Printf("% x\n", b)

	fmt.Printf("% x\n", big.NewInt(10).Bytes())

	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		fmt.Println("listen err", err)
		return
	}
	defer l.Close()

	n := new(netio)
	e := new(encode)
	n.Encode = e

	go n.handle_conn(ch_conn)
	if n.Encode != nil {
		fmt.Println("encode of netio is not nil")
	} else {
		fmt.Println("encode of netio is nil")
	}

	n.Encode.pack(nil)

	n.Encode = *new(Encode)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept err")
			return
		}

		ch_conn <- c
	}
}

func (net *netio) handle_conn(ch_conn chan net.Conn) {
	for conn := range ch_conn {
		go echo(conn)

		switch t := conn.(type) {
		default:
			fmt.Printf("unexcepted type %T\n", t)
		case io.Writer:
			fmt.Println("io.writer")
			//			fallthrough
		case Encode:
			fmt.Println("bool")

		}
	}
}

func echo(conn net.Conn) {
	if conn == nil {
		return
	}
	defer conn.Close()

	for {
		_, err := io.Copy(conn, conn)
		if err != nil {
			return
		}
	}
}

func pack(msg *Msg_package) (buf []byte) {

	return
}

func unpack(buf []byte) (msg *Msg_package) {

	return
}

func GetEncode_default() (en *Encode) {
	//	en = new(Encode{pack: pack, unpack:unpack})
	return
}
