package ws

import "time"

type Option struct {
	ReadBufferSize   int
	WriteBufferSize  int
	HandshakeTimeout time.Duration
}

const (
	DefaultReadBufferSize   = 1000
	DefaultWriteBufferSize  = 1000
	DefaultHandshakeTimeout = 10 * time.Second
)

func parseOption(option Option) Option {
	if option.ReadBufferSize == 0 {
		option.ReadBufferSize = DefaultReadBufferSize
	}

	if option.WriteBufferSize == 0 {
		option.WriteBufferSize = DefaultWriteBufferSize
	}

	if option.HandshakeTimeout == 0 {
		option.HandshakeTimeout = DefaultHandshakeTimeout
	}

	return option
}
