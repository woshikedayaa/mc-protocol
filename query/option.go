package query

import "time"

type Option func(client *BaseClient)
type OptionsType struct{}

func (op Option) apply(c *BaseClient) {
	op(c)
}

var Options OptionsType
var defaultOptions = []Option{Options.WithTimeOut(10 * time.Second)}

func (OptionsType) WithTimeOut(duration time.Duration) Option {
	return func(client *BaseClient) {
		client.timeout = duration
	}
}
