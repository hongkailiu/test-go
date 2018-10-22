package http

import "github.com/prometheus/client_golang/prometheus"

var (
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"device"},
	)
)
