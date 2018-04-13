package gobreak

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type runFunc func(context.Context) error
type fallbackFunc func(context.Context, error) error

// Do runs your function in synchronous manner, blocking until either your function succeeds
// or an error is returned, including circuit errors
func Do(ctx context.Context, name string, run runFunc, fall fallbackFunc) error {
	return <-Go(ctx, name, run, fall)
}

// Go runs your function in asynchronous manner, and returns error chan to the caller
func Go(ctx context.Context, name string, run runFunc, fall fallbackFunc) chan error {
	cmd := &command{
		name: name,
		// obtain circuit by name
		circuit: getCircuit(name),
		errChan: make(chan error, 1),
		run:     run,
		fall:    fall,
	}

	// ask circuit allow run or not
	done, err := cmd.circuit.Allow()
	if err != nil {
		cmd.errorWithFallback(ctx, err)
		return cmd.errChan
	}

	now := time.Now()
	once := sync.Once{}
	finished := make(chan struct{}, 0)

	go func() {
		// try recover when run function panics
		defer func() {
			if e := recover(); e != nil {
				once.Do(func() {
					done(false)
					cmd.elapsed = time.Now().Sub(now)
					cmd.errorWithFallback(ctx, fmt.Errorf("%s", e))
				})
			}
			// notify another goroutine
			finished <- struct{}{}
		}()

		// process run function
		err = run(ctx)

		// report run results to circuit
		once.Do(func() {
			done(err == nil)
			cmd.elapsed = time.Now().Sub(now)
			cmd.errorWithFallback(ctx, err)
		})
	}()

	// check if timeout or error happens
	go func() {
		select {
		case <-finished:
		case <-ctx.Done():
			once.Do(func() {
				done(false)
				cmd.elapsed = time.Now().Sub(now)
				cmd.errorWithFallback(ctx, ctx.Err())
			})
		}
	}()

	return cmd.errChan
}
