package query

import (
	"fmt"
	"testing"
	"time"
)

func TestBasicResponse_Encode(t *testing.T) {
	resp := &BasicResponse{}
	start := time.Now()
	err := resp.Encode(getBasicResponseTestExample())
	fmt.Println(time.Now().Sub(start).Milliseconds(), "ms")
	if err != nil {
		t.Fatal(err)
	}
	dumpBasic(resp)
}

func BenchmarkBasicResponse_Encode(b *testing.B) {
	resp := &BasicResponse{}
	bs := getBasicResponseTestExample()
	for i := 0; i < b.N; i++ {
		resp.Encode(bs)
	}
}

func getBasicResponseTestExample() []byte {
	/*
		Field name	Field Type			Example
		Type		byte				00
		Session-ID	int32				00 00 00 01
		MOTD		string				"A Minecraft Server\0"
		gametype	string				"SMP\0"
		map			string				"world\0"
		numplayers	string				"2\0"
		maxplayers	string				"20\0"
		hostport	Little-endian-short	DD 63 ( = 25565)
		hostip		string				"127.127.0.1\0"
	*/
	return []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x41, 0x20, 0x4D, 0x69, 0x6E, 0x65, 0x63, 0x72, 0x61, 0x66, 0x74, 0x20, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x00, 0x53, 0x4D, 0x50, 0x00, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x00, 0x32, 0x00, 0x32, 0x30, 0x00, 0xDD, 0x63, 0x31, 0x32, 0x37, 0x2E, 0x31, 0x32, 0x37, 0x2E, 0x30, 0x2E, 0x31, 0x00}
}

func dumpBasic(resp *BasicResponse) {
	fmt.Println("queryType", resp.EmptyResponse.typ)
	fmt.Println("sessionID", resp.EmptyResponse.sessionID)
	fmt.Println("MOTD", resp._MOTD)
	fmt.Println("gameType", resp.gameType)
	fmt.Println("map", resp._map)
	fmt.Println("numPlayer", resp.curPlayers)
	fmt.Println("maxPlayer", resp.maxPlayer)
	fmt.Println("port", resp.port)
	fmt.Println("ip", resp.ip)
}
