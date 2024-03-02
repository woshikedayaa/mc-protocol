package query

import (
	"net"
	"testing"
)

func TestUDPConnect(t *testing.T) {
	conn, err := net.Dial("udp", "172.30.28.83:5001")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("Hello\n"))
	if err != nil {
		t.Fatal(err)
	}
}
