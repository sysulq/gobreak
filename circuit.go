package gobreak

import (
	"sync"

	"github.com/sony/gobreaker"
)

var (
	circuitBreaker      map[string]*gobreaker.TwoStepCircuitBreaker
	circuitBreakerMutex *sync.RWMutex
)

func init() {
	circuitBreaker = make(map[string]*gobreaker.TwoStepCircuitBreaker)
	circuitBreakerMutex = &sync.RWMutex{}
}

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

func newCircuitBreaker(name string) *gobreaker.TwoStepCircuitBreaker {
	return gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
		Name: name,
	})
}
