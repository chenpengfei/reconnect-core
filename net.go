package reconnect_core

import (
	"context"
	"net"
)

type NetReconnection struct {
	net.Conn
	*Reconnection
}

func NewNetReconnection(ctx context.Context, network, address string, opts ...Option) *NetReconnection {
	nrc := &NetReconnection{
		Reconnection: NewReconnection(ctx, opts...),
	}

	nrc.operation = func() error {
		conn, err := net.Dial(network, address)
		if err == nil {
			nrc.Conn = conn
		}
		return err
	}

	nrc.Reconnection.retry()

	return nrc
}

func (nrc *NetReconnection) Close() error {
	err := nrc.Conn.Close()
	nrc.retry()
	return err
}
