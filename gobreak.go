package gobreak

import "time"

type runFunc func() error
type fallbackFunc func(error) error

const (
	Success   = "success"
	ErrReject = "reject"
	ErrFail   = "fail"
)

func Do(name string, run runFunc, fall fallbackFunc) error {
	c := getCircuit(name)

	done, err := c.Allow()
	if err != nil {
		requests.WithLabelValues(name, ErrReject).Inc()
		if fall != nil {
			err = fall(err)
		}
		return err
	}

	now := time.Now()

	err = run()

	elapsed := time.Now().Sub(now).Seconds()
	requestLatencyHistogram.WithLabelValues(name).Observe(elapsed)

	if err != nil {
		done(false)
		requests.WithLabelValues(name, ErrFail).Inc()
		if fall != nil {
			err = fall(err)
		}
		return err
	}

	done(true)
	requests.WithLabelValues(name, Success).Inc()
	return nil
}
