package reconnect_core

import (
	"context"
	"net"
)

type NetReconnection struct {
	net.Conn
	*Reconnection

	onConnect func(*NetReconnection)
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

	nrc.Reconnection.retry(nrc.retryDone)

	return nrc
}

func (nrc *NetReconnection) retryDone(err error) {
	if err == nil {
		nrc.onConnect(nrc)
	} else {
		nrc.onError(err)
	}
}

func (nrc *NetReconnection) Close() error {
	err := nrc.Conn.Close()
	nrc.retry(nrc.retryDone)
	return err
}

func (nrc *NetReconnection) OnConnect(onConnect func(*NetReconnection)) {
	nrc.onConnect = onConnect
}
