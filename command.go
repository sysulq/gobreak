package gobreak

import (
	"context"
	"errors"
	"time"

	"github.com/sony/gobreaker"
)

type command struct {
	name    string
	circuit *gobreaker.TwoStepCircuitBreaker
	errChan chan error
	run     runFunc
	fall    fallbackFunc
	elapsed time.Duration
}

var errPanic = errors.New("command panics")

// errorWithFallback process error and fallback logic, with prometheus metrics
func (c *command) errorWithFallback(ctx context.Context, err error) {

	// collect prometheus metrics
	event := "failure"
	switch err {
	case nil:
		event = "success"
	case context.DeadlineExceeded:
		event = "context-deadline-exceeded"
	case context.Canceled:
		event = "context-cancled"
	case gobreaker.ErrTooManyRequests:
		event = "too-many-requests"
	case gobreaker.ErrOpenState:
		event = "circuit-open"
	case errPanic:
		event = "panic"
	}

	requests.WithLabelValues(c.name, event).Inc()
	requestLatencyHistogram.WithLabelValues(c.name).Observe(c.elapsed.Seconds())

	// run returns nil means everything is ok
	if err == nil {
		c.errChan <- nil
		return
	}

	// return err directly when no fallback found
	if c.fall == nil {
		c.errChan <- err
		return
	}

	// fallback and return err
	err = c.fall(ctx, err)
	c.errChan <- err

	if err != nil {
		requests.WithLabelValues(c.name, "fallback-failure").Inc()
		return
	}
	requests.WithLabelValues(c.name, "fallback-success").Inc()
}
