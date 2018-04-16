package gobreak

import (
	"context"
	"fmt"
	"runtime"
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
		start:   time.Now(),
	}

	// ask circuit allow run or not
	done, err := cmd.circuit.Allow()
	if err != nil {
		cmd.errorWithFallback(ctx, err)
		return cmd.errChan
	}

	// Shared by the following two goroutines. It ensures only the faster
	// goroutine runs errorWithFallback().
	once := sync.Once{}
	finished := make(chan struct{}, 1)

	// goroutine 1
	go func() {
		// try recover when run function panics
		defer func() {
			// notify goroutine 2
			finished <- struct{}{}

			if e := recover(); e != nil {
				once.Do(func() {
					done(false)
					cmd.errorWithFallback(ctx, errPanic)

					stack := make([]byte, 1024*8)
					stack = stack[:runtime.Stack(stack, false)]
					fmt.Println(string(stack))
				})
			}
		}()

		// process run function
		err = run(ctx)

		// report run results to circuit
		once.Do(func() {
			done(err == nil)
			cmd.errorWithFallback(ctx, err)
		})
	}()

	// goroutine 2
	go func() {
		select {
		// check if goroutine 1 finished, timeout or error happens
		case <-finished:
		case <-ctx.Done():
			once.Do(func() {
				done(false)
				cmd.errorWithFallback(ctx, ctx.Err())
			})
		}
	}()

	return cmd.errChan
}
