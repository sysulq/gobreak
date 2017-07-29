package gobreak

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace = "gobreak"

	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "requests",
		Help:      "request count.",
	}, []string{"name", "state"})

	requestLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "request_latency_histogram",
		Help:      "request latency histogram.",
	}, []string{"name"})
)

func init() {
	prometheus.MustRegister(requests, requestLatencyHistogram)
}
