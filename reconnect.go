package reconnect_core

import (
	"context"
	"github.com/chenpengfei/backoff"
	"net"
)

type OnConnect func(conn *Reconnection)

type OnError func(err error)

type Strategy string

const (
	Fibonacci   Strategy = "fibonacci"
	Exponential Strategy = "exponential"
)

type Reconnection struct {
	net.Conn

	ctx      context.Context
	network  string
	address  string
	strategy Strategy

	onNotify  backoff.Notify
	onConnect OnConnect
	onError   OnError
}

func (re *Reconnection) OnConnect(onConnect OnConnect) {
	re.onConnect = onConnect
}

func (re *Reconnection) OnError(onError OnError) {
	re.onError = onError
}

func (re *Reconnection) OnNotify(onNotify backoff.Notify) {
	re.onNotify = onNotify
}

func (re *Reconnection) Close() error {
	err := re.Conn.Close()
	if err == nil {
		re.retry()
	}
	return err
}

func (re *Reconnection) retry() {
	go func() {
		var b backoff.BackOff
		switch re.strategy {
		case Fibonacci:
		case Exponential:
			b = backoff.NewExponentialBackOff()
		default:
			b = backoff.NewExponentialBackOff()
		}

		err := backoff.RetryNotify(
			func() error {
				conn, err := net.Dial(re.network, re.address)
				if err == nil {
					re.Conn = conn
				}
				return err
			},
			backoff.WithContext(b, re.ctx),
			re.onNotify)

		if err != nil {
			re.onError(err)
		} else {
			re.onConnect(re)
		}
	}()
}

func NewReconnection(address string, opts ...Option) *Reconnection {
	re := &Reconnection{
		ctx:      context.Background(),
		network:  "tcp",
		strategy: Exponential,
	}

	for _, opt := range opts {
		opt(re)
	}

	re.address = address

	re.onConnect = func(conn *Reconnection) {}
	re.onError = func(err error) {}

	re.retry()

	return re
}
