package main

import (
	"fmt"
	"github.com/woshikedayaa/mc-protocol/query"
)

func main() {
	client, err := query.NewQueryClient("127.0.0.1:5001")
	// with TimeOut. default 10s
	// client, err := query.NewQueryClient("127.0.0.1:5001",query.Options.WithTimeOut(10 * time.Second))
	if err != nil {
		panic(err)
	}
	defer client.Close()
	response, err := client.BasicRequest()
	if err != nil {
		panic(err)
	}
	// dump
	json, err := response.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
}
