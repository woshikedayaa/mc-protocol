package tcpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"testing"
)

func TestSimpleTcpClient_Read(t *testing.T) {
	addr := "127.0.0.1"
	port := 5001
	// go newTcpEchoServer(port)
	var client Client
	var err error
	client, err = NewIPV4Client(addr + ":" + strconv.Itoa(port))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Write(context.Background(), NewByteDataPack([]byte("Hello,world!")))
	if err != nil {
		t.Fatal(err)
	}
	var pack DataPack
	pack, err = client.Read()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(pack.String())
}

func TestNewTcpEchoServer(t *testing.T) {
	newTcpEchoServer(5001)
}

func newTcpEchoServer(port int) {
	listener, err := net.Listen("tcp4", "127.0.0.1"+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		log.Println("new client")
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			for {
				log.Println("read")
				readAll, err := io.ReadAll(c)
				log.Println("read finish")
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					log.Println(err)
					return
				}
				log.Println("content:", string(readAll))
				// 回写
				if len(readAll) != 0 {
					log.Println("write back")
					_, err := c.Write(readAll)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}(conn)
	}
}
