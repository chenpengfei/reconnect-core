package reconnect_core

import (
	"context"
	"github.com/chenpengfei/backoff"
	"net"
	"time"
)

type OnConnect func(rawConn net.Conn)

type OnError func(err error)

type OnNotify backoff.Notify

type Reconnection struct {
	backoff.BackOff

	onConnect OnConnect
	onNotify  OnNotify
	onError   OnError
}

func NewReconnection(ctx context.Context) *Reconnection {
	rc := &Reconnection{
		onConnect: func(conn net.Conn) {},
		onError:   func(err error) {},
		onNotify:  func(err error, duration time.Duration) {},
	}

	exp := backoff.NewExponentialBackOff()
	exp.MaxElapsedTime = backoff.Infinity
	rc.BackOff = backoff.WithContext(exp, ctx)

	return rc
}

func (rc *Reconnection) Dial(network, address string) {
	var raw net.Conn

	operation := func() error {
		conn, err := net.Dial(network, address)
		if err == nil {
			raw = conn
		}
		return err
	}

	err := backoff.RetryNotify(operation, rc.BackOff, backoff.Notify(rc.onNotify))
	if err != nil {
		rc.onError(err)
	} else {
		rc.onConnect(raw)
	}
}

func (rc *Reconnection) OnConnect(connect OnConnect) {
	rc.onConnect = connect
}

func (rc *Reconnection) OnError(onError OnError) {
	rc.onError = onError
}

func (rc *Reconnection) OnNotify(onNotify OnNotify) {
	rc.onNotify = onNotify
}
