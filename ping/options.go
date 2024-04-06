package ping

import (
	"errors"
	"reflect"
	"slices"
	"time"
)

type option func(client *Client)

func (o option) apply(c *Client) { o(c) }

type optionType struct {
	ops           []option
	timeout       time.Duration
	versionString string
	version       version
	network       string
}

func (ot *optionType) check(c *Client) error {
	ops := append(defaultOptions, ot.ops...)
	for i := 0; i < len(ot.ops); i++ {
		ops[i].apply(c)
	}
	var err, err2 error
	ot.version, err2 = newVersion(ot.versionString)
	if err2 != nil {
		err = errors.Join(err, err2)
	}
	return err // if ok, err == nil
}

var Options optionType

// for default options
var defaultOptions = []option{
	Options.WithNetwork("tcp"),
	Options.WithSpecialVersion("1.7"),
	Options.WithTimeout(10 * time.Second),
}

func (o option) String() string {
	return reflect.ValueOf(o).String()
}

func (ot *optionType) has(o option) bool {
	return slices.ContainsFunc(ot.ops, func(oi option) bool {
		return o.String() == oi.String()
	})
}

func (*optionType) WithSpecialVersion(s string) option {
	return func(c *Client) {
		c.op.versionString = s
	}
}

func (*optionType) WithTimeout(t time.Duration) option {
	return func(c *Client) {
		c.op.timeout = t
	}
}

func (*optionType) WithNetwork(network string) option {
	return func(c *Client) {
		c.op.network = network
	}
}
