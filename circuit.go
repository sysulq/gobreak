package gobreak

import (
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

var (
	// circuitBreaker holds a map filled with TwoStepCircuitBreaker
	circuitBreaker map[string]*gobreaker.TwoStepCircuitBreaker
	// circuitBreakerMutex is a RWMutex lock for circuitBreaker
	circuitBreakerMutex *sync.RWMutex
)

func init() {
	circuitBreaker = make(map[string]*gobreaker.TwoStepCircuitBreaker)
	circuitBreakerMutex = &sync.RWMutex{}
}

// getCircuit returns a TwoStepCircuitBreaker by name
func getCircuit(name string) *gobreaker.TwoStepCircuitBreaker {
	circuitBreakerMutex.RLock()
	cb, ok := circuitBreaker[name]
	if !ok {
		circuitBreakerMutex.RUnlock()
		circuitBreakerMutex.Lock()
		defer circuitBreakerMutex.Unlock()
		// because we released the rlock before we obtained the exclusive lock,
		// we need to double check that some other thread didn't beat us to
		// creation.
		if cb, ok := circuitBreaker[name]; ok {
			return cb
		}
		cb = newCircuitBreaker(name)
		circuitBreaker[name] = cb
	} else {
		defer circuitBreakerMutex.RUnlock()
	}

	return cb
}

// newCircuitBreaker creates TwoStepCircuitBreaker with suitable settings.
//
// Name is the name of the CircuitBreaker.
//
// MaxRequests is the maximum number of requests allowed to pass through
// when the CircuitBreaker is half-open.
// If MaxRequests is 0, the CircuitBreaker allows only 1 request.
//
// Interval is the cyclic period of the closed state
// for the CircuitBreaker to clear the internal Counts.
// If Interval is 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
//
// Timeout is the period of the open state,
// after which the state of the CircuitBreaker becomes half-open.
// If Timeout is 0, the timeout value of the CircuitBreaker is set to 60 seconds.
//
// ReadyToTrip is called with a copy of Counts whenever a request fails in the closed state.
// If ReadyToTrip returns true, the CircuitBreaker will be placed into the open state.
// If ReadyToTrip is nil, default ReadyToTrip is used.
//
// Default settings:
// MaxRequests: 3
// Interval:    5 * time.Second
// Timeout:     10 * time.Second
// ReadyToTrip: DefaultReadyToTrip
func newCircuitBreaker(name string) *gobreaker.TwoStepCircuitBreaker {
	return gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: DefaultReadyToTrip,
	})
}

// DefaultReadyToTrip returns true when the number of consecutive failures is more than 3 and rate of failure is more than 60%.
func DefaultReadyToTrip(counts gobreaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.ConsecutiveFailures >= 3 && failureRatio >= 0.6
}
