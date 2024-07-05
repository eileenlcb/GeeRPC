package geerpc

import (
	"net"
	"testing"
	"time"
)

type Bar int

func (b Bar) Timeout(argv int, reply *int) error {
	time.Sleep(2 * time.Second)
	return nil
}

func startServer(addr chan string) {
	var b Bar
	_ = Register(&b)

	l, _ := net.Listen("tcp", ":0")
	addr <- l.Addr().String()
	Accept(l)
}

func TestClient_dialTimeout(t *testing.T) {
	t.Parallel()
}
