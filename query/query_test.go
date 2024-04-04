package query

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 127.0.0.1 is author's local loop addr ,same as 127.0.0.1

func TestBaseClient_New(t *testing.T) {
	client, err := NewQueryClient("127.0.0.1:5001")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	response, err := client.HandShakeRequest()
	if err != nil {
		t.Fatal(err)
	}
	j, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(j))
}

func TestBaseClient_FullRequest(t *testing.T) {
	client, err := NewQueryClient("127.0.0.1:5001")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	response, err := client.FullRequest()
	if err != nil {
		t.Fatal(err)
	}
	j, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(j))
}

func TestBaseClient_BasicRequest(t *testing.T) {
	client, err := NewQueryClient("127.0.0.1:5001")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	response, err := client.BasicRequest()
	if err != nil {
		t.Fatal(err)
	}
	j, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(j))
}

func BenchmarkBaseClient_BasicRequest(b *testing.B) {
	client, err := NewQueryClient("127.0.0.1:5001")
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
	client, err := NewQueryClient("127.0.0.1:5001")
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
