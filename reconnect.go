package reconnect_core

import (
	"context"
	"github.com/chenpengfei/backoff"
	"time"
)

type OnError func(err error)

type Strategy string

const (
	Fibonacci   Strategy = "fibonacci"
	Exponential Strategy = "exponential"
)

type Reconnection struct {
	operation backoff.Operation
	backoff   backoff.BackOff

	strategy            Strategy
	initialInterval     time.Duration
	randomizationFactor float64
	multiplier          float64
	maxInterval         time.Duration
	maxElapsedTime      time.Duration

	onNotify backoff.Notify
	onError  OnError

	retrying bool
}

func NewReconnection(ctx context.Context, opts ...Option) *Reconnection {
	re := &Reconnection{
		strategy:            Exponential,
		initialInterval:     backoff.DefaultInitialInterval,
		randomizationFactor: backoff.DefaultRandomizationFactor,
		multiplier:          backoff.DefaultMultiplier,
		maxInterval:         backoff.DefaultMaxInterval,
		maxElapsedTime:      backoff.DefaultMaxElapsedTime,
		onError:             func(err error) {},
		onNotify:            func(err error, duration time.Duration) {},
		retrying:            false,
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

	re.backoff = backoff.WithContext(b, ctx)

	return re
}

func (rc *Reconnection) retry(done func(error)) {
	if rc.retrying {
		return
	}
	rc.retrying = true

	go func() {
		done(backoff.RetryNotify(rc.operation, rc.backoff, rc.onNotify))
		rc.retrying = false
	}()
}

func (rc *Reconnection) OnError(onError OnError) {
	rc.onError = onError
}

func (rc *Reconnection) OnNotify(onNotify backoff.Notify) {
	rc.onNotify = onNotify
}
