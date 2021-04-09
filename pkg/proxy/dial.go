package proxy

import (
	"context"
	corenet "net"
)

type DialContextFunc func(ctx context.Context, network, addr string) (corenet.Conn, error)
