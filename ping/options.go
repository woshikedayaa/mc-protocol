package ping

import (
	"errors"
	"github.com/woshikedayaa/mc-protocol/internal/ver"
	"reflect"
	"slices"
	"time"
)

type option func(client *Client)

func (o option) apply(c *Client) { o(c) }

type options struct {
	ops           []option
	timeout       time.Duration
	versionString string
	version       ver.Version
	network       string
}

func (ot *options) check(c *Client) error {
	ops := append(defaultOptions, ot.ops...)
	for i := 0; i < len(ops); i++ {
		ops[i].apply(c)
	}
	var err, err2 error
	ot.version, err2 = ver.ParseVersion(ot.versionString)
	if err2 != nil {
		err = errors.Join(err, err2)
	}
	return err // if ok, err == nil
}

var Options options

// for default options
var defaultOptions = []option{
	Options.WithNetwork("tcp"),
	Options.WithSpecialVersion("1.7"),
	Options.WithTimeout(10 * time.Second),
}

func (o option) String() string {
	return reflect.ValueOf(o).String()
}

func (ot *options) has(o option) bool {
	return slices.ContainsFunc(ot.ops, func(oi option) bool {
		return o.String() == oi.String()
	})
}

func (*options) WithSpecialVersion(s string) option {
	return func(c *Client) {
		c.op.versionString = s
	}
}

func (*options) WithTimeout(t time.Duration) option {
	return func(c *Client) {
		c.op.timeout = t
	}
}

func (*options) WithNetwork(network string) option {
	return func(c *Client) {
		c.op.network = network
	}
}
