package query

import "time"

type Option func(client *BaseClient)
type optionsType struct{}

func (op Option) apply(c *BaseClient) {
	op(c)
}

var Options optionsType

var defaultOptions = []Option{
	Options.WithTimeOut(10 * time.Second),
	Options.WithSpecialNetwork("udp"),
}

func (optionsType) WithTimeOut(duration time.Duration) Option {
	return func(client *BaseClient) {
		client.timeout = duration
	}
}

func (optionsType) WithSpecialNetwork(network string) Option {
	return func(client *BaseClient) {
		client.network = network
	}
}
