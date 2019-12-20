package reconnect_core

import (
	"time"
)

type Option func(*Reconnection)

func WithStrategy(strategy Strategy) Option {
	return func(reconnection *Reconnection) {
		reconnection.strategy = strategy
	}
}

func WithInitialInterval(interval time.Duration) Option {
	return func(reconnection *Reconnection) {
		reconnection.initialInterval = interval
	}
}

func WithRandomizationFactor(factor float64) Option {
	return func(reconnection *Reconnection) {
		reconnection.randomizationFactor = factor
	}
}

func WithMultiplier(multiplier float64) Option {
	return func(reconnection *Reconnection) {
		reconnection.multiplier = multiplier
	}
}

func WithMaxInterval(interval time.Duration) Option {
	return func(reconnection *Reconnection) {
		reconnection.maxInterval = interval
	}
}

func WithMaxElapsedTime(maxElapsedTime time.Duration) Option {
	return func(reconnection *Reconnection) {
		reconnection.maxElapsedTime = maxElapsedTime
	}
}
