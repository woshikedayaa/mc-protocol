package rcon

import (
	"fmt"
	"net"
	"testing"
)

func TestTcpRead(t *testing.T) {
	conn, err := net.Dial("tcp4", "172.30.28.83:5001")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("Hello\n"))
	if err != nil {
		t.Fatal(err)
	}
	var n int
	buf := make([]byte, 4)
	n, err = conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(buf[:n]))

	n, err = conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(buf[:n]))
}

func TestNewClient(t *testing.T) {
	client, err := NewClient("172.30.28.83:25575")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
}

func TestBaseClient_Auth(t *testing.T) {
	client, err := NewClient("172.30.28.83:25575")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	err = client.Auth("minecraft")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBaseClient_SendCommand(t *testing.T) {
	client, err := NewClient("172.30.28.83:25575")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	err = client.Auth("123456")
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.SendCommand("list")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(response.body))
	response, err = client.SendCommand("data get entity @e[limit=1]")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(response.body)
}
