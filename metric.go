package gobreak

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace = "gobreak"

	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "requests",
		Help:      "gobreak request count.",
	}, []string{"name", "state"})

	requestLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "request_latency_histogram",
		Help:      "gobreak request latency histogram.",
	}, []string{"name"})
)

func init() {
	prometheus.MustRegister(Collectors()...)
}

// Collectors returns all prometheus metric collectors
func Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		requests,
		requestLatencyHistogram,
	}
}
