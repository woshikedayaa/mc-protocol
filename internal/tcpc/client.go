package tcpc

import (
	"context"
	"io"
	"net"
)

type Client interface {
	// Connection Get the tcp connection
	Connection() net.Conn
	// Write a data-pack to target server
	Write(context.Context, DataPack) (int, error)
	// Read a data-pack from target server
	Read() (DataPack, error)
	// ReadN n must bigger than 0
	ReadN(int) (DataPack, error)
	// reader get io.Reader
	reader() io.Reader
	// Close the tcp
	Close() error
}
