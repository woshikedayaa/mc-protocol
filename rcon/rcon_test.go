package rcon

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewRconClient("127.0.0.1:5002")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
}

func TestBaseClient_Auth(t *testing.T) {
	client, err := NewRconClient("127.0.0.1:5002")
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
	client, err := NewRconClient("127.0.0.1:5002")
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
