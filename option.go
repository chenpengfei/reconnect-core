package reconnect_core

import (
	"context"
)

type Option func(*Reconnection)

func WithContext(ctx context.Context) Option {
	return func(reconnection *Reconnection) {
		reconnection.ctx = ctx
	}
}

func WithStrategy(strategy Strategy) Option {
	return func(reconnection *Reconnection) {
		reconnection.strategy = strategy
	}
}
