package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/hnlq715/gobreak"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		err := gobreak.Do("test", func() error {
			return errors.New("mock error\n")
		}, func(error) error {
			return errors.New("fallback\n")
		})
		rw.Write([]byte(err.Error()))
	})

	http.HandleFunc("/timeout", func(rw http.ResponseWriter, r *http.Request) {
		err := gobreak.Do("timeout", func() error {
			time.Sleep(2 * time.Second)
			return errors.New("mock error\n")
		}, nil)
		rw.Write([]byte(err.Error()))
	})

	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(os.Getpid(), ""))
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe("0.0.0.0:8000", nil)
}
