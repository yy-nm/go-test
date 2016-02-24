package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		return
	}
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)

	go func() {
		io.Copy(stdout, conn)
	}()

	go func() {
		io.Copy(conn, stdin)
	}()
	for {
		time.Sleep(time.Second * 2)
	}
}
