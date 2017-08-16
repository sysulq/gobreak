package gobreak

import "time"

type runFunc func() error
type fallbackFunc func(error) error

const (
	Success     = "success"
	ErrReject   = "reject"
	ErrFallBack = "fallback"
	ErrFail     = "fail"
	ErrPanic    = "panic"
)

// Do runs your function in a synchronous manner, blocking until either your function succeeds
// or an error is returned, including circuit errors
func Do(name string, run runFunc, fall fallbackFunc) error {
	errorType := Success
	// obtain circuit by name
	c := getCircuit(name)

	// ask circuit allow run or not
	done, err := c.Allow()
	if err != nil {
		errorType = ErrReject
		if fall != nil {
			errorType = ErrFallBack
			err = fall(err)
		}
		requests.WithLabelValues(name, errorType).Inc()
		return err
	}

	now := time.Now()

	// process run function
	err = run()
	// try recover when run function panics
	defer func() {
		e := recover()
		if e != nil {
			done(false)
			requests.WithLabelValues(name, ErrPanic).Inc()
			panic(e)
		}
	}()

	elapsed := time.Now().Sub(now).Seconds()
	requestLatencyHistogram.WithLabelValues(name).Observe(elapsed)

	// report run results to circuit
	done(err == nil)
	if err != nil {
		errorType = ErrFail
		if fall != nil {
			errorType = ErrFallBack
			err = fall(err)
		}
	}

	requests.WithLabelValues(name, errorType).Inc()
	return err
}
