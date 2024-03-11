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
	response, err := client.HandShakeRequest()
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

func BenchmarkBaseClient_BasicRequest(b *testing.B) {
	client, err := NewQueryClient("debian:5001")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()
	err = client.RefreshToken()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_, err = client.BasicRequest()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBaseClient_FullRequest(b *testing.B) {
	client, err := NewQueryClient("debian:5001")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()
	err = client.RefreshToken()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_, err = client.FullRequest()
		if err != nil {
			b.Fatal(err)
		}
	}
}
