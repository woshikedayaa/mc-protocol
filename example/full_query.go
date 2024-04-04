package main

import (
	"encoding/json"
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
	response, err := client.FullRequest()
	if err != nil {
		panic(err)
	}
	// dump
	j, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
