package app

import (
	"context"
	"crypto/tls"

	"github.com/mediocregopher/radix/v4"
)

func newPool(addr string) (radix.Client, error) {
	pool := radix.PoolConfig{}
	if isTLS {
		pool.Dialer.NetDialer = &tls.Dialer{
			NetDialer: nil,
			Config:    nil,
		}
	}

	return pool.New(context.Background(), "tcp", addr)
}
