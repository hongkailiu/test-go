package http

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed.",
		},
		[]string{"path", "hostname", "method"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "http request duration seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"path", "hostname", "method", "status"},
	)

	randomNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "random_number",
			Help: "the value of random number.",
		}, []string{"key", "hostname"},
	)

	//https://github.com/kubernetes/kubernetes/blob/master/pkg/volume/util/metrics.go
	storageOperationMetric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_operation_duration_seconds",
			Help:    "Storage operation duration",
			Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 15, 25, 50},
		},
		[]string{"volume_plugin", "operation_name", "hostname"},
	)
)

func prometheusRegister(component, selector string, client prowJobClient) *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registerProwJobCollector(component, client, selector, registry)
	registry.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	registry.MustRegister(httpRequestsTotal)
	registry.MustRegister(httpRequestDuration)
	registry.MustRegister(randomNumber)
	registry.MustRegister(storageOperationMetric)
	return registry
}
