# MC-protocol

Golang's implementation of modern MC related Procotol

## TODO

- [x] Query
- [x] Rcon
- [ ] Client

# Quick-start

```bash
# make sure your GO111MODULE=on
$ go get -u github.com/woshikedayaa/mc-protocol
...
$ curl https://raw.githubusercontent.com/woshikedayaa/mc-protocol/main/example/docker-compose.yaml \
> docker-compose.yaml
$ docker compose up -d # start minecraft-server
```

## BasicQuery

```go
// Query -> example/basic_query.go
func main() {
	client, err := query.NewQueryClient("127.0.0.1:5001")
	// with TimeOut. default 10s
	// client, err := query.NewQueryClient("127.0.0.1:5001",query.Options.WithTimeOut(10 * time.Second))
	if err != nil {
            panic(err)
	}
	defer client.Close() // optional
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
```

## Rcon

```go
// Rcon->example/rcon.go
func main() {
    client, err := rcon.NewRconClient("127.0.0.1:5002")
    // with TimeOut. default 10s
    // client, err := rcon.NewRconClient("debian:5001", rcon.Options.WithTimeOut(10 * time.Second))
    if err != nil {
    panic(err)
    }
    defer client.Close() // must
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
```

