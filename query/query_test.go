package query

import (
	"fmt"
	"testing"
)

func TestBaseClient_New(t *testing.T) {
	client, err := NewQueryClient("debian:5001")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	response, err := client.HandShake()
	if err != nil {
		t.Fatal(err)
	}
	json, err := response.JSON()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))
}

func TestBaseClient_FullRequest(t *testing.T) {
	client, err := NewQueryClient("debian:5001")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	response, err := client.FullRequest()
	if err != nil {
		t.Fatal(err)
	}
	json, err := response.JSON()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))
}

func TestBaseClient_BasicRequest(t *testing.T) {
	client, err := NewQueryClient("debian:5001")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	response, err := client.BasicRequest()
	if err != nil {
		t.Fatal(err)
	}
	json, err := response.JSON()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))
}
