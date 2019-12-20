package reconnect_core

import (
	"context"
	"github.com/chenpengfei/backoff"
	"io"
	"net"
	"time"
)

type OnConnect func(conn *Reconnection)

type OnError func(err error)

type Strategy string

const (
	Fibonacci   Strategy = "fibonacci"
	Exponential Strategy = "exponential"
)

type Reconnection struct {
	io.ReadWriteCloser

	operation backoff.Operation
	backoff   backoff.BackOff

	strategy            Strategy
	initialInterval     time.Duration
	randomizationFactor float64
	multiplier          float64
	maxInterval         time.Duration
	maxElapsedTime      time.Duration

	onNotify  backoff.Notify
	onConnect OnConnect
	onError   OnError

	retrying bool
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
	err := re.ReadWriteCloser.Close()
	re.retry()
	return err
}

func (re *Reconnection) retry() {
	if re.retrying {
		return
	}
	re.retrying = true

	go func() {
		err := backoff.RetryNotify(re.operation, re.backoff, re.onNotify)

		if err != nil {
			re.onError(err)
		} else {
			re.onConnect(re)
		}

		re.retrying = false
	}()
}

func NewReconnection(ctx context.Context, network, address string, opts ...Option) *Reconnection {
	re := &Reconnection{
		strategy:            Exponential,
		initialInterval:     backoff.DefaultInitialInterval,
		randomizationFactor: backoff.DefaultRandomizationFactor,
		multiplier:          backoff.DefaultMultiplier,
		maxInterval:         backoff.DefaultMaxInterval,
		maxElapsedTime:      backoff.DefaultMaxElapsedTime,
	}

	for _, opt := range opts {
		opt(re)
	}

	var initExponentialBackOff = func(eb *backoff.ExponentialBackOff) backoff.BackOff {
		eb.InitialInterval = re.initialInterval
		eb.RandomizationFactor = re.randomizationFactor
		eb.Multiplier = re.multiplier
		eb.MaxInterval = re.maxInterval
		eb.MaxElapsedTime = re.maxElapsedTime
		return eb
	}

	var b backoff.BackOff
	switch re.strategy {
	case Fibonacci:
	case Exponential:
		b = initExponentialBackOff(backoff.NewExponentialBackOff())
	default:
		b = initExponentialBackOff(backoff.NewExponentialBackOff())
	}

	re.operation = func() error {
		conn, err := net.Dial(network, address)
		if err == nil {
			re.ReadWriteCloser = conn
		}
		return err
	}
	re.backoff = backoff.WithContext(b, ctx)

	re.onConnect = func(conn *Reconnection) {}
	re.onError = func(err error) {}

	re.retry()

	return re
}
