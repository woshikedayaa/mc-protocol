package main

import (
	"fmt"
	"github.com/woshikedayaa/mc-protocol/rcon"
)

func main() {
	client, err := rcon.NewRconClient("127.0.0.1:5002")
	// // with TimeOut. default 10s
	// client, err := rcon.NewRconClient("debian:5001", rcon.Options.WithTimeOut(10*time.Second))
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.Auth("123456")
	if err != nil {
		panic(err)
	}
	response, err := client.SendCommand("list")
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}
